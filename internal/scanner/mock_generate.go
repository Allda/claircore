package scanner

//go:generate -command mockgen mockgen -package=scanner -self_package=github.com/quay/claircore/internal/scanner
//go:generate mockgen -destination=./scanner_mock.go github.com/quay/claircore/internal/scanner Scanner
//go:generate mockgen -destination=./store_mock.go github.com/quay/claircore/internal/scanner Store
//go:generate mockgen -destination=./fetcher_mock.go github.com/quay/claircore/internal/scanner Fetcher
//go:generate mockgen -destination=./layerscanner_mock.go github.com/quay/claircore/internal/scanner LayerScanner
//go:generate mockgen -destination=./packagescanner_mock.go github.com/quay/claircore/internal/scanner PackageScanner
//go:generate mockgen -destination=./versionedscanner_mock.go github.com/quay/claircore/internal/scanner VersionedScanner
