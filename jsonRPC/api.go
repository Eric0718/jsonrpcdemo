package jsonRPC

import (
	"errors"
	"fmt"
	"jsonrpcdemo/xwrap/util"
	"log"
)

//Returns the current chainId.
func (s *Server) eth_chainId() string {
	return uint64ToHexString(s.chainId)
}

//Returns the current network id.
func (s *Server) net_version() string {
	return s.networkId
}

//Returns the current client version.
func (s *Server) web3_clientVersion() string {
	return s.clinetVersion
}

//send signed transaction
func (s *Server) Eth_sendRawTransaction(rawTx string) (string, error) {
	log.Println("Into eth_sendRawTransaction===========", rawTx)

	etx, err := DecodeRawTx(rawTx)
	if err != nil {
		return "", err
	}
	//log.Printf("ethtx params:{to:%v,amount:%v,nounce:%v,hash:%v,gas:%v,gasPrice:%v,txType:%v,chainID:%v,tx lenght:%v}\n", etx.To(), etx.Value(), etx.Nonce(), etx.Hash(), etx.Gas(), etx.GasPrice(), etx.Type(), etx.ChainId(), len(etx.Data()))

	//check chainId
	if etx.ChainId().Uint64() != s.chainId {
		return "", fmt.Errorf("Wrong chainId,expect %v,got:%v", s.chainId, etx.ChainId().Int64())
	}

	//verify eth tx
	err = VerifyEthSignature(etx)
	if err != nil {
		return "", err
	}

	//convert eth tx to Top tx
	if !util.WrapEthTx(rawTx) {
		return "", errors.New("sendRawTransaction failed!")
	}

	return etx.Hash().Hex(), nil
}
