#!/bin/sh

echo "request chainId start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' localhost:37389

echo -e "\n\n request net wrok id start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' localhost:37389

echo -e "\n\n request web3_clientVersion start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":1}' localhost:37389

echo -e "\n\n send raw transaction start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0xf8920d85174876e80082520894d6139ea5fe0f3b54499e771417b0a5f56cd629b7880de0b6b3a7640000a477fb2c640000000000000000000000000000000000000000000000000de0b6b3a76400008240dea068374558f2dba5934f525aaf840a4e04d0506a33f94c5491f44db976f5f023f2a072caad5814801defb6c5fa3b0e7e6740fa264233673bd78912b11f439aa37aa9"],"id":1}' localhost:37389

echo -e "\n\n eth_blockNumber start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}' localhost:37389

echo -e "\n\n eth_getCode start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getCode","params":["0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b", "0x2"],"id":1}'  localhost:37389

echo -e "\n\n eth_getTransactionCount start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0x407d73d8a49eeb85d32cf465507dd71d507100c1","latest"],"id":1}' localhost:37389

echo -e "\n\n eth_gasPrice start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":73}' localhost:37389 

echo -e "\n\n eth_getBalance start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x407d73d8a49eeb85d32cf465507dd71d507100c1", "latest"],"id":1}' localhost:37389

echo -e "\n\n eth_call start>>>>>>>>>>"
curl -X POST --data '{"id":337305,"jsonrpc":"2.0","method":"eth_call","params":[{"data":"0x06fdde03","from":"0x0000000000000000000000000000000000000000","to":"0xd6139ea5fe0f3b54499e771417b0a5f56cd629b7"},"latest"]}' localhost:37389

echo -e "\n\n eth_estimateGas start>>>>>>>>>>"
curl -X POST --data '{"id":337305,"jsonrpc":"2.0","method":"eth_estimateGas","params":[{"data":"0x06fdde03","from":"0x0000000000000000000000000000000000000000","to":"0xd6139ea5fe0f3b54499e771417b0a5f56cd629b7"},"latest"]}' localhost:37389

echo -e "\n\n eth_getBlockByHash start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByHash","params":["0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae", false],"id":1}' localhost:37389

echo -e "\n\n eth_getBlockByNumber start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1b4", true],"id":1}' localhost:37389

echo -e "\n\n eth_getTransactionByHash start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"],"id":1}' localhost:37389

echo -e "\n\n eth_getTransactionReceipt start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"],"id":1}' localhost:37389

echo -e "\n\n eth_getStorageAt start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0", "method": "eth_getStorageAt", "params": ["0x295a70b2de5e3953354a6a8344e616ed314d7251", "0x0", "latest"], "id": 1}' localhost:37389

echo -e "\n\n eth_getLogs start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getLogs","params":[{"topics":["0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b"]}],"id":74}' localhost:37389

echo -e "\n\n web3_sha3 start>>>>>>>>>>"
curl -X POST --data '{"jsonrpc":"2.0","method":"web3_sha3","params":["0x68656c6c6f20776f726c64"],"id":64}' localhost:37389