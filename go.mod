module github.com/frankffenn/worker-srv

go 1.14

require (
	github.com/filecoin-project/filecoin-ffi v0.30.3
	github.com/filecoin-project/go-bitfield v0.0.2-0.20200518150651-562fdb554b6e // indirect
	github.com/filecoin-project/specs-actors v0.7.1
	github.com/filecoin-project/specs-storage v0.1.0
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/smartystreets/goconvey v1.6.4 // indirect
	go.uber.org/zap v1.14.1 // indirect
	golang.org/x/tools v0.0.0-20200108195415-316d2f248479 // indirect
)

replace github.com/filecoin-project/filecoin-ffi => ../filecoin-ffi

replace github.com/filecoin-project/sector-storage => ../my-sector-storage
