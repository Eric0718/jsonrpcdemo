go build -buildmode=c-shared -o libjsonrpc.so main.go
mv libjsonrpc.h /home/lyle/.local/include
mv libjsonrpc.so /home/lyle/.local/lib

gcc -o ./bin/runServer ./xcgo/main.c -ljsonrpc 
