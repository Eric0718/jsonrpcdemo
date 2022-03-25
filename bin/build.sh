go build -buildmode=c-shared -o libjsonrpc.so ../jsonrpc/run/main.go
mv libjsonrpc.* ./lib

export LD_LIBRARY_PATH=./lib:$LD_LIBRARY_PATH

gcc -o ./runServer ../xwrap/main.c -I./lib -L./lib -ljsonrpc 


