package jsonrpcserver

import (
	"encoding/binary"
	"encoding/json"
	"jsonrpcdemo/jsonrpc/client"
	"jsonrpcdemo/jsonrpc/util"
	"strconv"

	"fmt"
	"io/ioutil"
	"log"

	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

//New Server
func NewJsonRpcServer(chainId, networkId, archivePoint, clinetVersion string) *Server {
	cid, err := strconv.ParseUint(chainId[2:], 16, 32)
	if err != nil {
		panic(err)
	}
	return &Server{client: client.NewClient(), chainId: cid, networkId: networkId, archivePoint: archivePoint, clinetVersion: clinetVersion}
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
			REST = util.ResponseErrFunc(ParameterErr, "2.0", 0, err.(error).Error())
		}
	}()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("ioutil ReadAll error:%v\n", err)
		REST = util.ResponseErrFunc(IoutilErr, "", 0, err.Error())
		return
	}

	reqData := make(map[string]interface{})
	if err := json.Unmarshal(body, &reqData); err != nil {
		REST = util.ResponseErrFunc(JsonUnmarshalErr, "", 0, err.Error())
		return
	}

	method, err := util.GetString(reqData, "method")
	if err != nil {
		log.Printf("get method error:%v\n", err)
		REST = util.ResponseErrFunc(ParameterErr, "", 0, err.Error())
		return
	}

	jsonrpc, err := util.GetString(reqData, "jsonrpc")
	if err != nil {
		log.Printf("get jsonrpc error:%v\n", err)
		REST = util.ResponseErrFunc(ParameterErr, "", 0, err.Error())
		return
	}

	id, err := util.GetValue(reqData, "id")
	if err != nil {
		log.Printf("getValue error:%v\n", err)
		REST = util.ResponseErrFunc(ParameterErr, jsonrpc, 0, err.Error())
		return
	}

	switch method {
	case ETH_CHAINID:
		resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: s.eth_chainId()})
		if err != nil {
			log.Println("eth_chainId Marshal error:", err)
			REST = util.ResponseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
		} else {
			fmt.Println("eth_chainId success res>>>", s.eth_chainId())
			REST = resp
		}
	case NET_VERSION:
		resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: s.net_version()})
		if err != nil {
			log.Println("net_version Marshal error:", err)
			REST = util.ResponseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
		} else {
			log.Println("net_version success res>>>", s.net_version())
			REST = resp
		}
	case WEB3_CLIENTVERSION:
		resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: s.web3_clientVersion()})
		if err != nil {
			log.Println("web3_clientVersion Marshal error:", err)
			REST = util.ResponseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
		} else {
			fmt.Println("web3_clientVersion success :", s.web3_clientVersion())
			REST = resp
		}
	case ETH_SENDRAWTRANSACTION:
		para, err := util.GetRaw(reqData)
		if err != nil {
			log.Println("getRaw error:", err)
			REST = util.ResponseErrFunc(ParameterErr, jsonrpc, id, err.Error())
		} else {
			hash, err := s.Eth_sendRawTransaction(para[0].(string))
			if err != nil {
				log.Println("eth_sendRawTransaction error:", err)
				REST = util.ResponseErrFunc(UnkonwnErr, jsonrpc, id, err.Error())
			} else {
				resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: hash})
				if err != nil {
					log.Println("eth_sendRawTransaction Marshal error:", err)
					REST = util.ResponseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
				} else {
					log.Println("eth_sendRawTransaction success hash>>>", hash)
					REST = resp
				}
			}
		}
	case ETH_CALL:
		ret, err := s.Eth_call(reqData)
		if err != nil {
			log.Println("eth_call error:", err)
			var RetErr util.ErrorBody
			RetErr.Code = -4677
			RetErr.Message = err.Error()
			if len(ret) > 0 {
				btret := common.Hex2Bytes(ret)
				lenth := binary.BigEndian.Uint32(btret[64:68])
				data := btret[68 : lenth+68]
				errMsg := string(data)
				RetErr.Message = RetErr.Message + ": " + errMsg
				RetErr.Data = util.StringToHex(ret)
			}

			resp, err := json.Marshal(util.ResponseErr{JsonRPC: jsonrpc, Id: id, Error: &RetErr})
			if err != nil {
				log.Println("eth_call Marshal error:", err)
				resE := util.ResponseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
				w.Write(resE)
			} else {
				w.Write(resp)
			}
		} else {
			res := util.StringToHex(ret)
			log.Println("eth_call success res>>>", res)
			resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: res})
			if err != nil {
				log.Println("eth_call Marshal error:", err)
				resE := util.ResponseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
				w.Write(resE)
			} else {
				w.Write(resp)
				log.Println("return eth_call res length>>>", len(res))
			}
		}
	case ETH_BLOCKNUMBER:
		num, err := s.Eth_blockNumber()
		if err != nil {
			resE := util.ResponseErrFunc(UnkonwnErr, jsonrpc, id, err.Error())
			w.Write(resE)
		} else {
			resNum := fmt.Sprintf("%X", num)
			resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: (util.StringToHex(resNum))})
			if err != nil {
				log.Println("eth_blockNumber Marshal error:", err)
				resE := util.ResponseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
				w.Write(resE)
			} else {
				log.Println("eth_blockNumber success res>>>", util.StringToHex(resNum))
				w.Write(resp)
			}
		}
	case ETH_GETBALANCE:
		para, err := util.GetParam(reqData)
		if err != nil || len(para) == 0 {
			log.Println("util.GetParam error:", err)
			resE := util.ResponseErrFunc(ParameterErr, jsonrpc, id, err.Error())
			w.Write(resE)
		} else {
			from := para[0].(string)
			blc, err := s.Eth_getBalance(from)
			if err != nil {
				log.Println("eth_getBalance error:", err)
				resE := util.ResponseErrFunc(UnkonwnErr, jsonrpc, id, err.Error())
				w.Write(resE)
			} else {
				//metamask's decimal is 18,kto is 11,we need do blc*Pow10(7).
				// bigB := new(big.Int).SetUint64(blc)
				// bl := bigB.Mul(bigB, big.NewInt(ETHKTODIC))

				// resBalance := fmt.Sprintf("%X", bl)

				resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: blc})
				if err != nil {
					log.Println("eth_getBalance Marshal error:", err)
					resE := util.ResponseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
					w.Write(resE)
				} else {
					log.Println("eth_getBalance success res>>>", from, util.StringToHex(blc))
					w.Write(resp)
				}
			}
		}
	case ETH_GASPRICE:
		price, err := s.Eth_gasPrice()
		if err != nil {
			resE := util.ResponseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
			w.Write(resE)
		}
		resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: price})
		if err != nil {
			log.Println("eth_gasPrice Marshal error:", err)
			resE := util.ResponseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
			w.Write(resE)
		} else {
			w.Write(resp)
		}
	case EHT_GETCODE:
		para, err := util.GetParam(reqData)
		if err != nil || len(para) == 0 {
			log.Println("util.GetParam error:", err)
			resE := util.ResponseErrFunc(ParameterErr, jsonrpc, id, err.Error())
			w.Write(resE)
		} else {
			code, err := s.Eth_getCode(para[0].(string))
			if err != nil {
				log.Println("eth_getCode error:", err)
				code = "0x"
			}

			resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: code})
			if err != nil {
				log.Println("eth_getCode Marshal error:", err)
				resE := util.ResponseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
				w.Write(resE)
			} else {
				log.Println("eth_getCode success res>>>", code)
				w.Write(resp)
			}
		}
	case ETH_GETTRANSACTIONCOUNT:
		para, err := util.GetParam(reqData)
		if err != nil || len(para) == 0 {
			log.Println("util.GetParam error:", err)
			resE := util.ResponseErrFunc(ParameterErr, jsonrpc, id, err.Error())
			w.Write(resE)
		} else {
			addr := para[0].(string)
			count, err := s.Eth_getTransactionCount(addr)
			if err != nil {
				log.Println("eth_getTransactionCount error:", err)
				resE := util.ResponseErrFunc(UnkonwnErr, jsonrpc, id, err.Error())
				w.Write(resE)
			} else {
				//hexCount := fmt.Sprintf("%X", count)
				resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: count})
				if err != nil {
					log.Println("eth_getTransactionCount Marshal error:", err)
					resE := util.ResponseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
					w.Write(resE)
				} else {
					log.Println("eth_getTransactionCount success res>>>", "addr:", addr, "nonce", count)
					w.Write(resp)
				}
			}
		}
	case ETH_ESTIMATEGAS:
		ret, err := s.Eth_estimateGas(reqData)
		if err != nil {
			var RetErr util.ErrorBody
			RetErr.Code = -4677
			RetErr.Message = err.Error()
			if len(ret) > 0 {
				btret := common.Hex2Bytes(ret)
				lenth := binary.BigEndian.Uint32(btret[64:68])
				data := btret[68 : lenth+68]
				errMsg := string(data)
				RetErr.Message = RetErr.Message + ": " + errMsg
				RetErr.Data = util.StringToHex(ret)
			}

			resp, err := json.Marshal(util.ResponseErr{JsonRPC: jsonrpc, Id: id, Error: &RetErr})
			if err != nil {
				log.Println("eth_estimateGas Marshal error:", err)
				resE := util.ResponseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
				w.Write(resE)
			} else {
				log.Println("eth_estimateGas success ret>>>", ret)
				w.Write(resp)
			}
		} else {
			res := util.StringToHex(ret)
			log.Println("eth_estimateGas success res>>>", res)
			resp, err := json.Marshal(responseBody{JsonRPC: jsonrpc, Id: id, Result: res})
			if err != nil {
				log.Println("eth_estimateGas Marshal error:", err)
				resE := util.ResponseErrFunc(JsonMarshalErr, jsonrpc, id, err.Error())
				w.Write(resE)
			} else {
				log.Println("eth_estimateGas success res>>>", res)
				w.Write(resp)
			}
		}
	default:
		log.Printf("Error unsupport method:%v\n", method)
		REST = util.ResponseErrFunc(UnkonwnErr, jsonrpc, id, fmt.Errorf("Unsupport method:%v", method).Error())
	}
	log.Println("end HandRequest >>>>>>>>>>>>>>>")
}
