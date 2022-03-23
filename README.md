How to run：
    1、打开一个终端： 
        #./runServer
    2、打开另一个终端执行： 
        //request chainId
        curl -X POST --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' localhost:37389
        //request net wrok id
        curl -X POST --data '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' localhost:37389
        //request web3_clientVersion
        curl -X POST --data '{"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":1}' localhost:37389
        //send raw transaction
        curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0xf8920d85174876e80082520894d6139ea5fe0f3b54499e771417b0a5f56cd629b7880de0b6b3a7640000a477fb2c640000000000000000000000000000000000000000000000000de0b6b3a76400008240dea068374558f2dba5934f525aaf840a4e04d0506a33f94c5491f44db976f5f023f2a072caad5814801defb6c5fa3b0e7e6740fa264233673bd78912b11f439aa37aa9"],"id":1}' localhost:37389

