go build -buildmode=c-shared -o libjsonrpc.so ../main.go
mv libjsonrpc.* ./lib

gcc -o ./runServer ../xwrap/main.c -I./lib -L./lib -ljsonrpc 
