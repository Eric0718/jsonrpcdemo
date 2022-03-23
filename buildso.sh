go build -buildmode=c-shared -o libjsonrpc.so main.go
mv libjsonrpc.h ./xcgo
mv libjsonrpc.so ./xcgo/lib

gcc -o runServer ./xcgo/main.c ./xcgo/lib/libjsonrpc.so -I./xcgo
