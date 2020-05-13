module github.com/allda/claircore

go 1.12

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200227195959-4d242818bf55
	github.com/docker/docker => github.com/docker/docker v1.4.2-0.20200227233006-38f52c9fec82
)

require (
	github.com/asottile/dockerfile v3.1.0+incompatible
	github.com/buildkite/interpolate v0.0.0-20181028012610-973457fa2b4c
	github.com/crgimenes/goconfig v1.2.1
	github.com/docker/docker v1.13.1 // indirect
	github.com/doug-martin/goqu/v8 v8.6.0
	github.com/golang/mock v1.3.1
	github.com/google/go-cmp v0.4.0
	github.com/google/go-containerregistry v0.0.0-20191206185556-eb7c14b719c6
	github.com/google/uuid v1.1.1
	github.com/jackc/pgtype v1.0.0
	github.com/jackc/pgx/v4 v4.0.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/klauspost/compress v1.9.4
	github.com/knqyf263/go-cpe v0.0.0-20180327054844-659663f6eca2
	github.com/knqyf263/go-deb-version v0.0.0-20190517075300-09fca494f03d
	github.com/knqyf263/go-rpm-version v0.0.0-20170716094938-74609b86c936
	github.com/moby/buildkit v0.7.1 // indirect
	github.com/quay/alas v1.0.1
	github.com/quay/goval-parser v0.7.0
	github.com/remind101/migrate v0.0.0-20170729031349-52c1edff7319
	github.com/rs/zerolog v1.15.0
	github.com/shopspring/decimal v0.0.0-20190905144223-a36b5d85f337 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/tadasv/go-dpkg v0.0.0-20160704224136-c2cf9188b763
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/tools v0.0.0-20191205215504-7b8c8591a921
	gopkg.in/yaml.v3 v3.0.0-20191010095647-fc94e3f71652
)
