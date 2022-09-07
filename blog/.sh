brew services start mongodb-community@6.0
brew services stop mongodb-community@6.0
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
protoc --go_out=./blogpb --go-grpc_out=./blogpb --go-grpc_opt=require_unimplemented_servers=false ./blogpb/blog.proto