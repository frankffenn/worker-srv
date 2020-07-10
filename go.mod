module github.com/frankffenn/worker-srv

go 1.14

require (
	github.com/filecoin-project/filecoin-ffi v0.30.3
	github.com/filecoin-project/go-bitfield v0.0.2-0.20200518150651-562fdb554b6e // indirect
	github.com/filecoin-project/specs-actors v0.6.1
	github.com/filecoin-project/specs-storage v0.1.0
	github.com/go-kit/kit v0.8.0
	github.com/golang/protobuf v1.4.2
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/service/v2 v2.9.1 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	go.uber.org/zap v1.14.1 // indirect
	golang.org/x/tools v0.0.0-20200108195415-316d2f248479 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
)

replace github.com/coreos/etcd => github.com/ozonru/etcd v3.3.20-grpc1.27-origmodule+incompatible

replace github.com/filecoin-project/filecoin-ffi => ../filecoin-ffi

replace github.com/filecoin-project/sector-storage => ../my-sector-storage

//replace github.com/filecoin-project/specs-storage => ../my-specs-storage
