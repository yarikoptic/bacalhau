package streamingcid

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/filecoin-project/bacalhau/pkg/config"
	"github.com/filecoin-project/bacalhau/pkg/executor/docker"
	"github.com/filecoin-project/bacalhau/pkg/ipfs"
	"github.com/filecoin-project/bacalhau/pkg/model"
	"github.com/filecoin-project/bacalhau/pkg/storage"
	"github.com/filecoin-project/bacalhau/pkg/system"
	"github.com/rs/zerolog/log"
)

type StorageProvider struct {
	LocalDir   string
	IPFSClient ipfs.Client
	// XXX should be list of folders per channel so we can have multiple listeners
	// per channel per node. Also, cleanup obvs
	ChannelToFolderMap map[string]string
	// TODO mutex
}

// XXX why do we get instantiated multiples times motherfucker?

var TheOnlyStreamingCidStorageProvider *StorageProvider

func NewStorage(cm *system.CleanupManager, cl ipfs.Client) (*StorageProvider, error) {
	if TheOnlyStreamingCidStorageProvider == nil {
		dir, err := os.MkdirTemp(config.GetStoragePath(), "bacalhau-streaming-cid")
		if err != nil {
			return nil, err
		}

		cm.RegisterCallback(func() error {
			if err := os.RemoveAll(dir); err != nil {
				return fmt.Errorf("unable to clean up IPFS storage directory: %w", err)
			}
			return nil
		})

		storageHandler := &StorageProvider{
			IPFSClient:         cl,
			LocalDir:           dir,
			ChannelToFolderMap: map[string]string{},
		}

		err = storageHandler.setupStreamingGossipsub()
		if err != nil {
			return nil, err
		}

		log.Trace().Msgf("Streaming CID driver created with address: %s", cl.APIAddress())
		TheOnlyStreamingCidStorageProvider = storageHandler
		return storageHandler, nil
	} else {
		return TheOnlyStreamingCidStorageProvider, nil
	}
}

func (dockerIPFS *StorageProvider) IsInstalled(ctx context.Context) (bool, error) {
	_, err := dockerIPFS.IPFSClient.ID(ctx)
	return err == nil, err
}

func (dockerIPFS *StorageProvider) HasStorageLocally(ctx context.Context, volume model.StorageSpec) (bool, error) {
	return true, nil
}

// we wrap this in a timeout because if the CID is not present on the network this seems to hang
func (dockerIPFS *StorageProvider) GetVolumeSize(ctx context.Context, volume model.StorageSpec) (uint64, error) {
	return 0, nil
}

func (dockerIPFS *StorageProvider) PrepareStorage(ctx context.Context, storageSpec model.StorageSpec) (storage.StorageVolume, error) {
	ctx, span := system.GetTracer().Start(ctx, "storage/ipfs/apicopy.PrepareStorage")
	defer span.End()

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, config.GetDownloadCidRequestTimeout(ctx))
	defer cancel()

	channelLocalFolder := path.Join(dockerIPFS.LocalDir, storageSpec.Channel)
	// make a new folder from the Source property of the storage volume
	// appended onto LocalDir - don't worry if the folder already exists
	err := os.MkdirAll(channelLocalFolder, os.ModePerm)
	if err != nil {
		return storage.StorageVolume{}, err
	}

	log.Printf(">>>>>>>>> Adding to channelToFolderMap[%s] = %s", storageSpec.Channel, channelLocalFolder)
	dockerIPFS.ChannelToFolderMap[storageSpec.Channel] = channelLocalFolder

	return storage.StorageVolume{
		Type:   storage.StorageVolumeConnectorBind,
		Source: channelLocalFolder,
		Target: storageSpec.Path,
	}, nil
}

//nolint:lll // Exception to the long rule
func (dockerIPFS *StorageProvider) CleanupStorage(_ context.Context, storageSpec model.StorageSpec, _ storage.StorageVolume) error {
	return nil
}

func (dockerIPFS *StorageProvider) Upload(ctx context.Context, localPath string) (model.StorageSpec, error) {
	return model.StorageSpec{}, nil
}

func (dockerIPFS *StorageProvider) Explode(ctx context.Context, spec model.StorageSpec) ([]model.StorageSpec, error) {
	return []model.StorageSpec{}, nil
}

var DID_SUBSCRIBE bool

func (dockerIPFS *StorageProvider) setupStreamingGossipsub() error {
	if docker.GlobalStreamingResultPubSubConnection == nil {
		return fmt.Errorf("GlobalStreamingResultPubSubConnection has not been created")
	}

	if DID_SUBSCRIBE {
		return nil
	}
	// GlobalStreamingResultPubSubConnection.Publish(ctx, model.StreamingResult{})
	err := docker.GlobalStreamingResultPubSubConnection.Subscribe(context.Background(), dockerIPFS)
	DID_SUBSCRIBE = true
	return err
}

func (dockerIPFS *StorageProvider) Handle(ctx context.Context, message model.CIDStreamElement) error {
	fmt.Printf("message inside storage driver --------------------------------------\n")
	spew.Dump(message)
	localFolderPath, ok := dockerIPFS.ChannelToFolderMap[message.Channel]
	if !ok {
		return fmt.Errorf("Streaming CID storage driver could not find a local folder for channel %s, have %+v", message.Channel, dockerIPFS.ChannelToFolderMap)
	}
	return dockerIPFS.IPFSClient.Get(ctx, message.CID, path.Join(localFolderPath, message.CID))
}

func (dockerIPFS *StorageProvider) getFileFromIPFS(ctx context.Context, storageSpec model.StorageSpec) (storage.StorageVolume, error) {
	ctx, span := system.GetTracer().Start(ctx, "storage/ipfs/apicopy.copyFile")
	defer span.End()

	outputPath := filepath.Join(dockerIPFS.LocalDir, storageSpec.CID)

	// If the output path already exists, we already have the data, as
	// ipfsClient.Get(...) renames the result path atomically after it has
	// finished downloading the CID.
	ok, err := system.PathExists(outputPath)
	if err != nil {
		return storage.StorageVolume{}, err
	}
	if !ok {
		err = dockerIPFS.IPFSClient.Get(ctx, storageSpec.CID, outputPath)
		if err != nil {
			return storage.StorageVolume{}, err
		}
	}

	volume := storage.StorageVolume{
		Type:   storage.StorageVolumeConnectorBind,
		Source: outputPath,
		Target: storageSpec.Path,
	}

	return volume, nil
}

// Compile time interface check:
var _ storage.Storage = (*StorageProvider)(nil)
