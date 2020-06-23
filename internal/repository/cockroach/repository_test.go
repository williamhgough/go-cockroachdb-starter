package cockroach_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"

	"github.com/williamhgough/go-cockroachdb-starter/internal/repository/cockroach"
)

var dbURL string

// set up using: https://github.com/ory/dockertest/blob/v3/examples/CockroachDB.md
func TestMain(m *testing.M) {
	var db *sql.DB
	var err error

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{Repository: "cockroachdb/cockroach", Tag: "v19.2.0", Cmd: []string{"start", "--insecure"}})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	dbURL = fmt.Sprintf("postgresql://root@localhost:%s/defaultdb?sslmode=disable", resource.GetPort("26257/tcp"))
	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to cockroach container: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestAddX(t *testing.T) {
	_, err := cockroach.NewRepository(dbURL)
	require.NoError(t, err)
}
