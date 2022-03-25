package jsonrpcserver

import (
	"encoding/json"
	"jsonrpcdemo/jsonrpc/util"
	"strconv"

	"fmt"
	"io/ioutil"
	"log"

	"net/http"
)

//New Server
func NewJsonRpcServer(chainId, networkId, archivePoint, clinetVersion string) *Server {
	cid, err := strconv.ParseUint(chainId[2:], 16, 32)
	if err != nil {
		panic(err)
	}
	return &Server{chainId: cid, networkId: networkId, archivePoint: archivePoint, clinetVersion: clinetVersion}
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
	default:
		log.Printf("Error unsupport method:%v\n", method)
		REST = util.ResponseErrFunc(UnkonwnErr, jsonrpc, id, fmt.Errorf("Unsupport method:%v", method).Error())
	}
	log.Println("end HandRequest >>>>>>>>>>>>>>>")
}
