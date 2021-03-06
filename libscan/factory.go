package libscan

import (
	"fmt"

	"github.com/quay/claircore/dpkg"
	"github.com/quay/claircore/internal/scanner"
	"github.com/quay/claircore/internal/scanner/defaultfetcher"
	"github.com/quay/claircore/internal/scanner/defaultlayerscanner"
	"github.com/quay/claircore/internal/scanner/defaultscanner"
	"github.com/quay/claircore/pkg/distlock"
	dlpg "github.com/quay/claircore/pkg/distlock/postgres"
)

// ScannerFactory is a factory method to return a Scanner interface during libscan runtime.
type ScannerFactory func(lib *libscan, opts *Opts) (scanner.Scanner, error)

// scannerFactory is the default ScannerFactory
func scannerFactory(lib *libscan, opts *Opts) (scanner.Scanner, error) {
	// add other distributed locking implementations here as they grow
	var sc distlock.Locker
	switch opts.ScanLock {
	case PostgresSL:
		sc = dlpg.NewLock(lib.db, opts.ScanLockRetry)
	default:
		return nil, fmt.Errorf("provided ScanLock opt is unsupported")
	}

	// add other fetcher implementations here as they grow
	var ft scanner.Fetcher
	ft = defaultfetcher.New(lib.client, nil, opts.LayerFetchOpt)

	sOpts := &scanner.Opts{
		Store:           lib.store,
		ScanLock:        sc,
		Fetcher:         ft,
		PackageScanners: lib.packageScanners(),
	}

	// add other layer scanner implementations as they grow
	sOpts.LayerScanner = defaultlayerscanner.New(opts.LayerScanConcurrency, sOpts)
	s := defaultscanner.New(sOpts)
	return s, nil
}

// PackageScannerFactory is a factory method to return a set of PackageScanners which are used during libscan runtime.
type PackageScannerFactory func() []scanner.PackageScanner

// packageScannerFactory is the default PackageScannerFactory
func packageScannerFactory() []scanner.PackageScanner {
	return []scanner.PackageScanner{
		// add other packge scanners as they grow
		dpkg.NewPackageScanner(),
	}
}
