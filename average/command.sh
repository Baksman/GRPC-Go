
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
protoc --go_out=./averagepb --go-grpc_out=./averagepb --go-grpc_opt=require_unimplemented_servers=false ./averagepb/average.proto