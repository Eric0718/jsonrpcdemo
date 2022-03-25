go build -buildmode=c-shared -o libjsonrpc.so ../jsonrpc/run/main.go
mv libjsonrpc.* ./lib

go build -buildmode=c-shared -o libgrpc.so ../grpc/run/main.go
mv libgrpc.* ./lib

#gcc -o ./runJsonRpc ../cgotest/main.c -I./lib -L./lib -ljsonrpc     #C run json rpc
gcc -o ./runGrpc ../cgotest/main.c -I./lib -L./lib -lgrpc      #C run grpc

