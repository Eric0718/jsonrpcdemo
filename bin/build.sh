#//build jsonrpc
go build -buildmode=c-shared -o libjsonrpc.so ../jsonrpc/run/main.go
mv libjsonrpc.* ./lib
gcc -o ./runJsonRpc ../jsonrpc/run/main.c -I./lib -L./lib -ljsonrpc     

#//build grpc
go build -buildmode=c-shared -o libgrpc.so ../grpc/run/main.go
mv libgrpc.* ./lib
gcc -o ./runGrpc ../grpc/run/main.c -I./lib -L./lib -lgrpc      #C run grpc

