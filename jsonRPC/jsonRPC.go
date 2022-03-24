package jsonRPC

import (
	"encoding/json"
	"errors"
	"math/big"
	"strconv"

	"fmt"
	"io/ioutil"
	"jsonrpcdemo/xcgo/util"
	"log"

	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

//New Server
func NewServer(chainId, networkId, consensusPoint, archivePoint, clinetVersion string) *Server {
	cid, err := strconv.ParseUint(chainId[2:], 16, 32)
	if err != nil {
		panic(err)
	}
	return &Server{chainId: cid, networkId: networkId, consensusPoint: consensusPoint, archivePoint: archivePoint, clinetVersion: clinetVersion}
}

//Handle Request
func (s *Server) HandRequest(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")

	var REST []byte
	defer func() {
		w.Write(REST)
	}()

	defer func() {
		if err := recover(); err != nil {
			log.Println("Error:", err)
			REST = responseErrFunc(ParameterErr, "2.0", 0, err.(error).Error())
		}
	}()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("ioutil ReadAll error:%v\n", err)
		REST = responseErrFunc(IoutilErr, "", 0, err.Error())
		return
	}

	reqData := make(map[string]interface{})
	if err := json.Unmarshal(body, &reqData); err != nil {
		REST = responseErrFunc(JsonUnmarshalErr, "", 0, err.Error())
		return
	}

	method, err := getString(reqData, "method")
	if err != nil {
		log.Printf("get method error:%v\n", err)
		REST = responseErrFunc(ParameterErr, "", 0, err.Error())
		return
	}

	jsonrpc, err := getString(reqData, "jsonrpc")
	if err != nil {
		log.Printf("get jsonrpc error:%v\n", err)
		REST = responseErrFunc(ParameterErr, "", 0, err.Error())
		return
	}

	id, err := getValue(reqData, "id")
	if err != nil {
		log.Printf("getValue error:%v\n", err)
		REST = responseErrFunc(ParameterErr, jsonrpc, 0, err.Error())
		return
	}

	switch method {
	case ETH_CHAINID:
		resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: s.eth_chainId()})
		if err != nil {
			log.Println("eth_chainId Marshal error:", err)
			REST = responseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
		} else {
			fmt.Println("eth_chainId success res>>>", s.eth_chainId())
			REST = resp
		}
	case NET_VERSION:
		resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: s.net_version()})
		if err != nil {
			log.Println("net_version Marshal error:", err)
			REST = responseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
		} else {
			log.Println("net_version success res>>>", s.net_version())
			REST = resp
		}
	case WEB3_CLIENTVERSION:
		resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: s.web3_clientVersion()})
		if err != nil {
			log.Println("web3_clientVersion Marshal error:", err)
			REST = responseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
		} else {
			fmt.Println("web3_clientVersion success :", s.web3_clientVersion())
			REST = resp
		}
	case ETH_SENDRAWTRANSACTION:
		para, err := getRaw(reqData)
		if err != nil {
			log.Println("getRaw error:", err)
			REST = responseErrFunc(ParameterErr, jsonrpc, id, err.Error())
		} else {
			hash, err := s.Eth_sendRawTransaction(para[0].(string))
			if err != nil {
				log.Println("eth_sendRawTransaction error:", err)
				REST = responseErrFunc(UnkonwnErr, jsonrpc, id, err.Error())
			} else {
				resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: hash})
				if err != nil {
					log.Println("eth_sendRawTransaction Marshal error:", err)
					REST = responseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
				} else {
					log.Println("eth_sendRawTransaction success hash>>>", hash)
					REST = resp
				}
			}
		}
	default:
		log.Printf("Error unsupport method:%v\n", method)
		REST = responseErrFunc(UnkonwnErr, jsonrpc, id, fmt.Errorf("Unsupport method:%v", method).Error())
	}
	log.Println("end HandRequest >>>>>>>>>>>>>>>")
}

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
	if !util.ConvertEthTx(rawTx) {
		return "", errors.New("sendRawTransaction failed!")
	}

	return etx.Hash().Hex(), nil
}

//Decode eth Transaction Data
func DecodeRawTx(rawTx string) (*types.Transaction, error) {
	body, err := hexutil.Decode(rawTx)
	if err != nil {
		return nil, err
	}
	var etx types.Transaction
	err = rlp.DecodeBytes(body, &etx)
	if err != nil {
		return nil, err
	}

	return &etx, nil
}

//parse eth signature
func parseEthSignature(ethtx *types.Transaction) []byte {
	big8 := big.NewInt(8)
	v, r, s := ethtx.RawSignatureValues()
	v = new(big.Int).Sub(v, new(big.Int).Mul(ethtx.ChainId(), big.NewInt(2)))
	v.Sub(v, big8)

	var sign []byte
	sign = append(sign, r.Bytes()...)
	sign = append(sign, s.Bytes()...)
	sign = append(sign, byte(v.Uint64()-27))
	return sign
}

//Verify Eth Signature
func VerifyEthSignature(ethtx *types.Transaction) error {
	sign := parseEthSignature(ethtx)
	if len(sign) <= 64 {
		return fmt.Errorf("eth signature lenght error:%v", len(sign))
	}
	pub, err := crypto.Ecrecover(ethtx.Hash().Bytes(), sign)
	if err != nil {
		return err
	}
	if !crypto.VerifySignature(pub, ethtx.Hash().Bytes(), sign[:64]) {
		return fmt.Errorf("Verify Eth Signature failed!")
	}
	return nil
}
