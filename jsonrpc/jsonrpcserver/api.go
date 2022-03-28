package jsonrpcserver

import (
	"errors"
	"fmt"
	"jsonrpcdemo/jsonrpc/util"
	"jsonrpcdemo/jsonrpc/xwrap"

	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/goinggo/mapstructure"
)

//Returns the current chainId.
func (s *Server) eth_chainId() string {
	return util.Uint64ToHexString(s.chainId)
}

//Returns the current network id.
func (s *Server) net_version() string {
	return s.networkId
}

//Returns the current client version.
func (s *Server) web3_clientVersion() string {
	return s.clinetVersion
}

//Returns the number of most recent block.
func (s *Server) Eth_blockNumber() (string, error) {
	blockH, err := s.client.GetMaxBlockNumber()
	if err != nil {
		return "", err
	}
	return util.Uint64ToHexString(blockH), nil
}

//Returns code at a given address.
func (s *Server) Eth_getCode(addr string) (string, error) {
	return s.client.GetCode(addr)
}

//Returns the number of transactions sent from an address.
func (s *Server) Eth_getTransactionCount(addr string) (string, error) {
	nonce, err := s.client.GetNonce(addr)
	if err != nil {
		return "", err
	}
	return util.Uint64ToHexString(nonce), nil
}

//Returns the current price per gas in wei.
func (s *Server) Eth_gasPrice() (string, error) {
	price, err := s.client.Eth_gasPrice()
	if err != nil {
		return "", err
	}
	return util.Uint64ToHexString(price), nil
}

//Returns a list of addresses owned by client.
func (s *Server) Eth_accounts() []common.Address {
	var accounts []common.Address
	return accounts
}

//Returns the balance of the account of given address.
func (s *Server) Eth_getBalance(from string) (string, error) {
	balance, err := s.client.GetBalance(from)
	if err != nil {
		return "", err
	}
	return util.Uint64ToHexString(balance), nil
}

//send signed transaction
func (s *Server) Eth_sendRawTransaction(rawTx string) (string, error) {
	log.Println("Into eth_sendRawTransaction===========", rawTx)

	etx, err := util.DecodeRawTx(rawTx)
	if err != nil {
		return "", err
	}
	//log.Printf("ethtx params:{to:%v,amount:%v,nounce:%v,hash:%v,gas:%v,gasPrice:%v,txType:%v,chainID:%v,tx lenght:%v}\n", etx.To(), etx.Value(), etx.Nonce(), etx.Hash(), etx.Gas(), etx.GasPrice(), etx.Type(), etx.ChainId(), len(etx.Data()))

	//check chainId
	if etx.ChainId().Uint64() != s.chainId {
		return "", fmt.Errorf("Wrong chainId,expect %v,got:%v", s.chainId, etx.ChainId().Int64())
	}

	//verify eth tx
	err = util.VerifyEthSignature(etx)
	if err != nil {
		return "", err
	}

	//convert eth tx to Top tx
	if !xwrap.WrapEthTx(rawTx) {
		return "", errors.New("sendRawTransaction failed!")
	}

	return etx.Hash().Hex(), nil
}

//Executes a new message call immediately without creating a transaction on the block chain.
func (s *Server) Eth_call(mp map[string]interface{}) (string, error) {
	v, ok := mp["params"]
	if !ok {
		return "", errors.New(fmt.Sprintf("'%s' not exist", "params"))
	}

	if _, ok := v.([]interface{}); !ok {
		return "", errors.New("eth_call: params is wrong!")
	}

	Para := v.([]interface{})
	var para params
	if len(Para) > 0 {
		err := mapstructure.Decode(Para[0].(map[string]interface{}), &para)
		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New(fmt.Errorf("eth_call: Decode Para[%v] failed!", Para).Error())
	}

	ret, _, err := s.client.Eth_call(para.From, para.To, para.Data, util.Check0x(para.Value))
	if err != nil {
		return ret, err
	}
	return ret, nil
}

//Generates and returns an estimate of how much gas is necessary to allow the transaction to complete. The transaction will not be added to the blockchain.
//Note that the estimate may be significantly more than the amount of gas actually used by the transaction, for a variety of reasons including EVM mechanics and node performance.
func (s *Server) Eth_estimateGas(mp map[string]interface{}) (string, error) {
	v, ok := mp["params"]
	if !ok {
		return "", errors.New(fmt.Sprintf("'%s' not exist", "params"))
	}

	if _, ok := v.([]interface{}); !ok {
		return "", errors.New("eth_estimateGas: params is wrong!")
	}

	Para := v.([]interface{})
	var para params
	if len(Para) > 0 {
		err := mapstructure.Decode(Para[0].(map[string]interface{}), &para)
		if err != nil {
			return "", err
		}
	}

	log.Printf("eth_estimateGas params: from=%v,to=%v,gas=%v,gasPrice=%v,value=%v,data lenght=%v\n", para.From, para.To, para.Gas, para.GasPrice, para.Value, len(para.Data))

	limit, err := s.client.Get_gasLimit()
	if err != nil {
		return "", err
	}

	if len(para.Data) <= 0 {
		return fmt.Sprintf("%X", limit), nil
	}

	if len(para.To) <= 0 {
		return fmt.Sprintf("%X", limit), nil
	}

	ret, gas, err := s.client.Eth_call(para.From, para.To, para.Data, util.Check0x(para.Value))
	if err != nil {
		log.Println("eth_estimateGas :", ret, "error", err)
		return ret, err
	}

	log.Println("eth_estimateGas successfully,gas:", gas)
	return gas, nil
}
