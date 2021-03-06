// +build integration

package postgres

import (
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/quay/claircore"
	"github.com/quay/claircore/test"
	pgtest "github.com/quay/claircore/test/postgres"
	"github.com/stretchr/testify/assert"
)

// scanInfo is a helper struct for providing scanner information
// in test tables
type scnrInfo struct {
	name    string `integration:"name"`
	version string `integration:"version"`
	kind    string `integration:"kind"`
	id      int64  `integration:"id"`
}

func Test_IndexPackages_Success_Parallel(t *testing.T) {
	tt := []struct {
		// the name of this benchmark
		name string
		// number of packages to index.
		pkgs int
		// the layer that holds the discovered packages
		layer *claircore.Layer
	}{
		{
			name: "10 packages",
			pkgs: 10,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
		{
			name: "50 packages",
			pkgs: 50,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
		{
			name: "100 packages",
			pkgs: 100,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
		{
			name: "250 packages",
			pkgs: 250,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
		{
			name: "500 packages",
			pkgs: 500,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
		{
			name: "1000 packages",
			pkgs: 1000,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
		{
			name: "2000 packages",
			pkgs: 2000,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
		{
			name: "3000 packages",
			pkgs: 3000,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
	}

	db, store, teardown := NewTestStore(t)
	defer teardown()

	// we will create a subtest which blocks the teardown() until
	// all parallel tests are finished
	t.Run("blocking-group", func(t *testing.T) {
		for _, tab := range tt {
			table := tab
			t.Run(table.name, func(t *testing.T) {
				t.Parallel()
				// gen a scnr and insert
				vscnrs := test.GenUniqueScanners(1)
				err := pgtest.InsertUniqueScanners(db, vscnrs)

				// gen packages
				pkgs := test.GenUniquePackages(table.pkgs)

				// run the indexing
				err = store.IndexPackages(pkgs, table.layer, vscnrs[0])
				if err != nil {
					t.Fatalf("failed to index packages: %v", err)
				}

				assert.NoError(t, err)
				checkScanArtifact(t, db, pkgs, table.layer)
			})
		}
	})
}

func Test_IndexPackages_Success(t *testing.T) {
	tt := []struct {
		// the name of this benchmark
		name string
		// number of packages to index.
		pkgs int
		// the layer that holds the discovered packages
		layer *claircore.Layer
	}{
		{
			name: "10 packages",
			pkgs: 10,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
		{
			name: "50 packages",
			pkgs: 50,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
		{
			name: "100 packages",
			pkgs: 100,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
		{
			name: "250 packages",
			pkgs: 250,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
		{
			name: "500 packages",
			pkgs: 500,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
		{
			name: "1000 packages",
			pkgs: 1000,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
		{
			name: "2000 packages",
			pkgs: 2000,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
		{
			name: "3000 packages",
			pkgs: 3000,
			layer: &claircore.Layer{
				Hash: "test-layer-hash",
			},
		},
	}

	for _, table := range tt {
		t.Run(table.name, func(t *testing.T) {
			db, store, teardown := NewTestStore(t)
			defer teardown()

			// gen a scnr and insert
			vscnrs := test.GenUniqueScanners(1)
			err := pgtest.InsertUniqueScanners(db, vscnrs)

			// gen packages
			pkgs := test.GenUniquePackages(table.pkgs)

			// run the indexing
			err = store.IndexPackages(pkgs, table.layer, vscnrs[0])
			if err != nil {
				t.Fatalf("failed to index packages: %v", err)
			}

			assert.NoError(t, err)
			checkScanArtifact(t, db, pkgs, table.layer)
		})
	}

}

// checkScanArtifact confirms a scanartifact is created linking the layer, package/source/distribution entities from the layer, and scnr which discovered these.
// indirectly we test that dists and packages are indexed correctly by querying with their unique fields.
func checkScanArtifact(t *testing.T, db *sqlx.DB, expectedPkgs []*claircore.Package, layer *claircore.Layer) {
	// NOTE: we gen one scanner for this test with ID 0, this is hard coded into this check
	for _, pkg := range expectedPkgs {
		var distID sql.NullInt64
		err := db.Get(
			&distID,
			`SELECT id FROM dist 
			WHERE name = $1 
			AND version = $2 
			AND version_code_name = $3 
			AND version_id = $4 
			AND arch = $5`,
			pkg.Dist.Name,
			pkg.Dist.Version,
			pkg.Dist.VersionCodeName,
			pkg.Dist.VersionID,
			pkg.Dist.Arch,
		)
		if err != nil {
			t.Fatalf("failed to query for distribution %v: %v", pkg.Dist, err)
		}

		var pkgID sql.NullInt64
		err = db.Get(
			&pkgID,
			`SELECT id FROM package 
			WHERE name = $1 
			AND kind = $2 
			AND version = $3`,
			pkg.Name,
			pkg.Kind,
			pkg.Version,
		)
		if err != nil {
			t.Fatalf("received error selecting package id %s version %s", pkg.Name, pkg.Version)
		}

		var layer_hash, sakind string
		var package_id, dist_id, scanner_id sql.NullInt64
		row := db.QueryRowx(
			`SELECT layer_hash, kind, package_id, dist_id, scanner_id 
			FROM scanartifact 
			WHERE layer_hash = $1 
			AND package_id = $2 
			AND scanner_id = $3`,
			layer.Hash,
			pkgID,
			0,
		)

		err = row.Scan(&layer_hash, &sakind, &package_id, &dist_id, &scanner_id)
		if err != nil {
			if err == sql.ErrNoRows {
				t.Fatalf("failed to find scanartifact for pkg %v", pkg)
			}
			t.Fatalf("received error selecting scanartifact for pkg %v: %v", pkg, err)
		}

		assert.Equal(t, layer.Hash, layer_hash)
		assert.Equal(t, "package", sakind)
		assert.Equal(t, pkgID.Int64, package_id.Int64)
		assert.Equal(t, distID.Int64, dist_id.Int64)
		assert.Equal(t, int64(0), scanner_id.Int64)
	}
}
