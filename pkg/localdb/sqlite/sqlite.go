package sqlite

import (
	"context"
	"embed"
	"fmt"
	"time"

	sync "github.com/lukemarsden/golang-mutex-tracer"

	"database/sql"

	"github.com/filecoin-project/bacalhau/pkg/localdb"
	"github.com/filecoin-project/bacalhau/pkg/model"
	"github.com/filecoin-project/bacalhau/pkg/system"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

type SQLiteDatastore struct {
	mtx      sync.RWMutex
	filename string
	db       *sql.DB
}

func NewSQLiteDatastore(filename string) (*SQLiteDatastore, error) {
	db, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, err
	}
	datastore := &SQLiteDatastore{
		db:       db,
		filename: filename,
	}
	err = datastore.migrate()
	if err != nil {
		return nil, err
	}
	datastore.mtx.EnableTracerWithOpts(sync.Opts{
		Threshold: 10 * time.Millisecond,
		Id:        "SQLiteDatastore.mtx",
	})
	return datastore, nil
}

// Gets a job from the datastore.
//
// Errors:
//
//   - error-job-not-found        		  -- if the job is not found
func (d *SQLiteDatastore) GetJob(ctx context.Context, id string) (*model.Job, error) {
	//nolint:ineffassign,staticcheck
	ctx, span := system.GetTracer().Start(ctx, "pkg/localdb/sqlite/SQLiteDatastore.GetJob")
	defer span.End()

	d.mtx.RLock()
	defer d.mtx.RUnlock()
	return nil, nil
}

// Get Job Events from a job ID
//
// Errors:
//
//   - error-job-not-found        		  -- if the job is not found
func (d *SQLiteDatastore) GetJobEvents(ctx context.Context, id string) ([]model.JobEvent, error) {
	//nolint:ineffassign,staticcheck
	ctx, span := system.GetTracer().Start(ctx, "pkg/localdb/sqlite/SQLiteDatastore.GetJobEvents")
	defer span.End()

	return []model.JobEvent{}, nil
}

func (d *SQLiteDatastore) GetJobLocalEvents(ctx context.Context, id string) ([]model.JobLocalEvent, error) {
	//nolint:ineffassign,staticcheck
	ctx, span := system.GetTracer().Start(ctx, "pkg/localdb/sqlite/SQLiteDatastore.GetJobLocalEvents")
	defer span.End()
	return []model.JobLocalEvent{}, nil
}

func (d *SQLiteDatastore) GetJobs(ctx context.Context, query localdb.JobQuery) ([]*model.Job, error) {
	ctx, span := system.GetTracer().Start(ctx, "pkg/localdb/sqlite/SQLiteDatastore.GetJobs")
	defer span.End()

	d.mtx.RLock()
	defer d.mtx.RUnlock()
	return []*model.Job{}, nil
}

func (d *SQLiteDatastore) HasLocalEvent(ctx context.Context, jobID string, eventFilter localdb.LocalEventFilter) (bool, error) {
	return false, nil
}

func (d *SQLiteDatastore) AddJob(ctx context.Context, j *model.Job) error {
	//nolint:ineffassign,staticcheck
	ctx, span := system.GetTracer().Start(ctx, "pkg/localdb/sqlite/SQLiteDatastore.AddJob")
	defer span.End()
	return nil
}

func (d *SQLiteDatastore) AddEvent(ctx context.Context, jobID string, ev model.JobEvent) error {
	//nolint:ineffassign,staticcheck
	ctx, span := system.GetTracer().Start(ctx, "pkg/localdb/sqlite/SQLiteDatastore.AddEvent")
	defer span.End()
	return nil
}

func (d *SQLiteDatastore) AddLocalEvent(ctx context.Context, jobID string, ev model.JobLocalEvent) error {
	//nolint:ineffassign,staticcheck
	ctx, span := system.GetTracer().Start(ctx, "pkg/localdb/sqlite/SQLiteDatastore.AddLocalEvent")
	defer span.End()
	return nil
}

func (d *SQLiteDatastore) UpdateJobDeal(ctx context.Context, jobID string, deal model.Deal) error {
	//nolint:ineffassign,staticcheck
	ctx, span := system.GetTracer().Start(ctx, "pkg/localdb/sqlite/SQLiteDatastore.UpdateJobDeal")
	defer span.End()
	return nil
}

func (d *SQLiteDatastore) GetJobState(ctx context.Context, jobID string) (model.JobState, error) {
	//nolint:ineffassign,staticcheck
	ctx, span := system.GetTracer().Start(ctx, "pkg/localdb/sqlite/SQLiteDatastore.GetJobState")
	defer span.End()

	return model.JobState{}, nil
}

func (d *SQLiteDatastore) UpdateShardState(
	ctx context.Context,
	jobID, nodeID string,
	shardIndex int,
	update model.JobShardState,
) error {
	//nolint:ineffassign,staticcheck
	ctx, span := system.GetTracer().Start(ctx, "pkg/localdb/sqlite/SQLiteDatastore.UpdateShardState")
	defer span.End()
	return nil
}

// helper method to read a single job from memory. This is used by both GetJob and GetJobs.
// It is important that we don't attempt to acquire a lock inside this method to avoid deadlocks since
// the callers are expected to be holding a lock, and golang doesn't support reentrant locks.
func (d *SQLiteDatastore) getJob(id string) (*model.Job, error) {
	return nil, nil
}

//go:embed migrations/*.sql
var fs embed.FS

func (d *SQLiteDatastore) migrate() error {
	files, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	migrations, err := migrate.NewWithSourceInstance("iofs", files, fmt.Sprintf("sqlite://%s", d.filename))
	if err != nil {
		return err
	}
	return migrations.Up()
}

// Static check to ensure that Transport implements Transport:
var _ localdb.LocalDB = (*SQLiteDatastore)(nil)
