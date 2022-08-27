protoc  greet/greetpb/greet.proto --go_out=plugins=grpc:.
protoc  calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.

protoc -I calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:calculator

<------------- right command -------------------->
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
protoc --go_out=./calculatorpb --go-grpc_out=./calculatorpb --go-grpc_opt=require_unimplemented_servers=false ./calculatorpb/calculator.proto