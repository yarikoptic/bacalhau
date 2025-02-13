//go:build unit || !integration

package downloader

import (
	"context"
	"crypto/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bacalhau-project/bacalhau/pkg/util/closer"

	ipfs2 "github.com/bacalhau-project/bacalhau/pkg/downloader/ipfs"
	"github.com/bacalhau-project/bacalhau/pkg/ipfs"

	"github.com/bacalhau-project/bacalhau/pkg/logger"
	"github.com/bacalhau-project/bacalhau/pkg/model"
	"github.com/bacalhau-project/bacalhau/pkg/system"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestDownloaderSuite(t *testing.T) {
	suite.Run(t, new(DownloaderSuite))
}

type DownloaderSuite struct {
	suite.Suite
	cm               *system.CleanupManager
	client           ipfs.Client
	outputDir        string
	downloadSettings *model.DownloaderSettings
	downloadProvider DownloaderProvider
}

func (ds *DownloaderSuite) SetupSuite() {
	logger.ConfigureTestLogging(ds.T())
	system.InitConfigForTesting(ds.T())
}

// Before each test
func (ds *DownloaderSuite) SetupTest() {
	ds.cm = system.NewCleanupManager()
	ds.T().Cleanup(func() {
		ds.cm.Cleanup(context.Background())
	})

	ctx, cancel := context.WithCancel(context.Background())
	ds.T().Cleanup(cancel)

	node, err := ipfs.NewLocalNode(ctx, ds.cm, nil)
	require.NoError(ds.T(), err)

	ds.client = node.Client()

	swarm, err := node.SwarmAddresses()
	require.NoError(ds.T(), err)

	testOutputDir := ds.T().TempDir()
	ds.outputDir = testOutputDir

	ds.downloadSettings = &model.DownloaderSettings{
		Timeout:        model.DefaultIPFSTimeout,
		OutputDir:      testOutputDir,
		IPFSSwarmAddrs: strings.Join(swarm, ","),
	}

	ds.downloadProvider = model.NewMappedProvider(
		map[model.StorageSourceType]Downloader{
			model.StorageSourceIPFS: ipfs2.NewIPFSDownloader(ds.cm, ds.downloadSettings),
		},
	)
}

// Generate a file with random data.
func generateFile(path string) ([]byte, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer closer.CloseWithLogOnError("file", file)

	b := make([]byte, 128)
	_, err = rand.Read(b)
	if err != nil {
		return nil, err
	}

	_, err = file.Write(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Wraps generation of a set of output files that represent the output from a
// specific result, and saves them to IPFS.
//
// The passed setup func will be called with a temporary directory. Within the
// setup func, the user should make a number of calls to `mockFile` to generate
// files within the directory. At the end, the entire directory is saved to
// IPFS.
func mockOutput(ds *DownloaderSuite, setup func(string)) string {
	testDir := ds.T().TempDir()

	setup(testDir)

	cid, err := ds.client.Put(context.Background(), testDir)
	require.NoError(ds.T(), err)

	return cid
}

// Generates a test file at the given path filled with random data, ensuring
// that any parent directories for the file are also present.
func mockFile(ds *DownloaderSuite, path ...string) []byte {
	filePath := filepath.Join(path...)
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	require.NoError(ds.T(), err)

	contents, err := generateFile(filePath)
	require.NoError(ds.T(), err)
	return contents
}

// Requires that a file exists when the path is traversed downwards from the
// output directory.
func requireFileExists(ds *DownloaderSuite, path ...string) string {
	testPath := filepath.Join(ds.outputDir, filepath.Join(path...))
	require.FileExistsf(ds.T(), testPath, "File %s not present", testPath)

	return testPath
}

// Requires that a file exists with the specified contents when the path is
// traversed downwards from the output directory.
func requireFile(ds *DownloaderSuite, expected []byte, path ...string) {
	testPath := requireFileExists(ds, path...)

	contents, err := os.ReadFile(testPath)
	require.NoError(ds.T(), err)
	require.Equal(ds.T(), expected, contents)
}

func (ds *DownloaderSuite) TestNoExpectedResults() {
	err := DownloadResults(
		context.Background(),
		[]model.PublishedResult{},
		ds.downloadProvider,
		ds.downloadSettings,
	)
	require.NoError(ds.T(), err)
}

func (ds *DownloaderSuite) TestFullOutput() {
	var exitCode, stdout, stderr, hello, goodbye []byte
	cid := mockOutput(ds, func(dir string) {
		exitCode = mockFile(ds, dir, "exitCode")
		stdout = mockFile(ds, dir, model.DownloadFilenameStdout)
		stderr = mockFile(ds, dir, "stderr")
		hello = mockFile(ds, dir, "outputs", "hello.txt")
		goodbye = mockFile(ds, dir, "outputs", "goodbye.txt")
	})

	err := DownloadResults(
		context.Background(),
		[]model.PublishedResult{
			{
				NodeID: "testnode",
				Data: model.StorageSpec{
					StorageSource: model.StorageSourceIPFS,
					Name:          "result-0",
					CID:           cid,
				},
			},
		},
		ds.downloadProvider,
		ds.downloadSettings,
	)
	require.NoError(ds.T(), err)

	requireFile(ds, stdout, "stdout")
	requireFile(ds, stderr, "stderr")
	requireFile(ds, exitCode, "exitCode")
	requireFile(ds, goodbye, "outputs", "goodbye.txt")
	requireFile(ds, hello, "outputs", "hello.txt")
}

func (ds *DownloaderSuite) TestOutputWithNoStdFiles() {
	cid := mockOutput(ds, func(dir string) {
		mockFile(ds, dir, "outputs", "lonely.txt")
	})

	err := DownloadResults(
		context.Background(),
		[]model.PublishedResult{
			{
				NodeID: "testnode",
				Data: model.StorageSpec{
					StorageSource: model.StorageSourceIPFS,
					Name:          "result-0",
					CID:           cid,
				},
			},
		},
		ds.downloadProvider,
		ds.downloadSettings,
	)
	require.NoError(ds.T(), err)

	requireFileExists(ds, "outputs", "lonely.txt")
}

func (ds *DownloaderSuite) TestCustomVolumeNames() {
	cid := mockOutput(ds, func(s string) {
		mockFile(ds, s, "secrets", "private.pem")
	})

	err := DownloadResults(
		context.Background(),
		[]model.PublishedResult{
			{
				NodeID: "testnode",
				Data: model.StorageSpec{
					StorageSource: model.StorageSourceIPFS,
					Name:          "result-0",
					CID:           cid,
				},
			},
		},
		ds.downloadProvider,
		ds.downloadSettings,
	)
	require.NoError(ds.T(), err)

	requireFileExists(ds, "secrets", "private.pem")
}
