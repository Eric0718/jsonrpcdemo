#include<stdio.h>
#include"libjsonrpc.h"
#include<thread>
#include<chrono>
#include <string>
#include <unistd.h>
#include <iostream>

GoString buildGoString(const char* p, size_t n){
    return {p, static_cast<ptrdiff_t>(n)};
}

int main() {
    std::string chainId = "0x538";
    std::string networkId = "0x538";
    std::string archivePoint = "127.0.0.1:37399";
    std::string clientversion = "clientversion:hello there!";
    std::string jsonrpcPort = "0.0.0.0:37389";

    void *srv;
    srv = RunJsonRpc(buildGoString(chainId.c_str(),chainId.size()),  
                    buildGoString(networkId.c_str(),networkId.size()),
                    buildGoString(archivePoint.c_str(),archivePoint.size()),
                    buildGoString(clientversion.c_str(),clientversion.size()),
                    buildGoString(jsonrpcPort.c_str(),jsonrpcPort.size()));
 
    //stop server in seconds
    std::this_thread::sleep_for(std::chrono::seconds(1600));  

    int code;  
    code = StopJsonRpc(srv);
    if (code != 0){
        printf("stop json rpc failed!! code=%d\n",code);
    }
    
    return 0;
}