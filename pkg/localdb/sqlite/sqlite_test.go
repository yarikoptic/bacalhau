//go:build unit || !integration

package sqlite

import (
	"io/ioutil"
	"log"
	"testing"

	"database/sql"

	_ "github.com/filecoin-project/bacalhau/pkg/logger"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestSQLiteMigrations(t *testing.T) {
	datafile, err := ioutil.TempFile("", "sqlite-test-*.db")
	require.NoError(t, err)
	_, err = NewSQLiteDatastore(datafile.Name())
	require.NoError(t, err)

	db, err := sql.Open("sqlite", datafile.Name())
	require.NoError(t, err)
	_, err = db.Exec(`
insert into jobs (job_id) values (123);
`)
	require.NoError(t, err)
	var id int
	rows, err := db.Query(`
select job_id from jobs;
`)
	require.NoError(t, err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id)
	}
	err = rows.Err()
	require.NoError(t, err)
	require.Equal(t, id, 123)
}
