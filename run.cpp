#include<stdio.h>
#include"libjsonrpc.h"
#include<thread>
#include<chrono>

int main() {
    printf("test cmake ok!!");

    void *srv;
    srv = RunJsonRpc();
 
     std::this_thread::sleep_for(std::chrono::seconds(10));
    StopJsonRpc(srv);
    
    return 0;
}