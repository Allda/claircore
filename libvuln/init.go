package libvuln

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/quay/claircore/internal/updater"
	"github.com/quay/claircore/internal/vulnstore"
	"github.com/quay/claircore/internal/vulnstore/postgres"
	pglock "github.com/quay/claircore/pkg/distlock/postgres"
)

// initUpdaters provides initial burst control to not launch too many updaters at once.
// returns any errors on eC and returns a CaneclFunc on dC to stop all updaters
func initUpdaters(opts *Opts, db *sqlx.DB, store vulnstore.Updater, dC chan context.CancelFunc, eC chan error) {
	// just to be defensive
	err := opts.Parse()
	if err != nil {
		eC <- err
		return
	}

	controllers := map[string]*updater.Controller{}

	for _, u := range opts.Updaters {
		if _, ok := controllers[u.Name()]; ok {
			eC <- fmt.Errorf("duplicate updater found in UpdaterFactory. all names must be unique: %s", u.Name())
			return
		}
		controllers[u.Name()] = updater.New(&updater.Opts{
			Updater:       u,
			Store:         store,
			Name:          u.Name(),
			Interval:      opts.UpdateInterval,
			Lock:          pglock.NewLock(db, time.Duration(0)),
			UpdateOnStart: false,
		})
	}

	// limit initial concurrent updates
	cc := make(chan struct{}, opts.UpdaterInitConcurrency)

	var wg sync.WaitGroup
	wg.Add(len(controllers))
	for _, v := range controllers {
		cc <- struct{}{}
		vv := v
		go func() {
			updateTO, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			err := vv.Update(updateTO)
			if err != nil {
				eC <- fmt.Errorf("updater %s failed to update: %v", vv.Name, err)
			}
			wg.Done()
			cancel()
			<-cc
		}()
	}
	wg.Wait()
	close(eC)

	// start all updaters and return context
	ctx, cancel := context.WithCancel(context.Background())
	for _, v := range controllers {
		v.Start(ctx)
	}
	dC <- cancel
}

// initStore initializes a vulsntore and returns the underlying db object also
func initStore(opts *Opts) (*sqlx.DB, vulnstore.Store, error) {
	switch opts.DataStore {
	case Postgres:
		// we are going to use pgx for more control over connection pool and
		// and a cleaner api around bulk inserts
		connconfig, err := pgx.ParseConnectionString(opts.ConnString)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse ConnString: %v", err)
		}
		pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
			ConnConfig:     connconfig,
			MaxConnections: 30,
			AfterConnect:   nil,
			AcquireTimeout: 30 * time.Second,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create ConnPool: %v", err)
		}

		// setup sqlx
		db := stdlib.OpenDBFromPool(pool)
		sqlxDB := sqlx.NewDb(db, "pgx")

		store := postgres.NewVulnStore(sqlxDB, pool)
		return sqlxDB, store, nil
	default:
		return nil, nil, fmt.Errorf("provided unknown DataStore")
	}
}
