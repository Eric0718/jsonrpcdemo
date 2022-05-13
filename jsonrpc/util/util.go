package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func GetString(mp map[string]interface{}, k string) (string, error) {
	v, ok := mp[k]
	if !ok {
		return "", errors.New(fmt.Sprintf("'%s' not exist", k))
	}
	if s, ok := v.(string); ok {
		return s, nil
	}
	return "", errors.New(fmt.Sprintf("'%s' not string", k))
}

func GetValue(mp map[string]interface{}, k string) (interface{}, error) {
	v, ok := mp[k]
	if !ok {
		return 0, errors.New(fmt.Sprintf("'%s' not exist", k))
	}
	//log.Printf("value type %T,value:%v\n", v, v)
	return v, nil
}

func GetParam(mp map[string]interface{}) ([]interface{}, error) {
	v, ok := mp["params"]
	if !ok {
		return nil, errors.New(fmt.Sprintf("'%s' not exist", "params"))
	}
	if _, ok := v.([]interface{}); !ok {
		return nil, errors.New("params are wrong!")
	}
	return v.([]interface{}), nil
}

func GetRaw(mp map[string]interface{}) ([]interface{}, error) {
	v, ok := mp["params"]
	if !ok {
		return nil, errors.New(fmt.Sprintf("'%s' not exist", "params"))
	}
	if len(v.([]interface{})) != 1 {
		return nil, fmt.Errorf("Wrong raw data:%v", v.([]interface{}))
	}

	return v.([]interface{}), nil
}

func Check0x(s string) string {
	if len(s) > 2 && s[:2] == "0x" {
		return s[2:]
	}
	return s
}

func Uint64ToHexString(val uint64) string {
	return fmt.Sprintf("0x%X", val)
}

func HexToUint64(hxs string) (uint64, error) {
	if len(hxs) > 2 {
		if hxs[:2] == "0x" {
			hxs = hxs[2:]
		}
	}
	n, err := strconv.ParseUint(hxs, 16, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

type ErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type ResponseErr struct {
	JsonRPC string      `json:"jsonrpc"`
	Id      interface{} `json:"id"`
	Error   *ErrorBody  `json:"error"`
}

func ResponseErrFunc(code int, jsonRpc string, id interface{}, msg string) []byte {
	Err := &ErrorBody{Code: code, Message: msg}
	resp, err := json.Marshal(ResponseErr{JsonRPC: jsonRpc, Id: id, Error: Err})
	if err != nil {
		log.Println("eth_sendTransaction Marshal error:", err)
		return []byte(err.Error())
	}
	return resp
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

func zeroBytes() []byte {
	return []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
}

//parse eth signature
func parseEthSignature(ethtx *types.Transaction) []byte {
	big8 := big.NewInt(8)
	v, r, s := ethtx.RawSignatureValues()
	v = new(big.Int).Sub(v, new(big.Int).Mul(ethtx.ChainId(), big.NewInt(2)))
	v.Sub(v, big8)

	rBytes := r.Bytes()
	if n := len(rBytes); n < 32 {
		rBytes = append(zeroBytes()[:32-n], rBytes...)
	}
	sBytes := s.Bytes()
	if n := len(sBytes); n < 32 {
		sBytes = append(zeroBytes()[:32-n], sBytes...)
	}
	vBytes := byte(v.Uint64() - 27)

	var sign []byte
	sign = append(sign, rBytes...)
	sign = append(sign, sBytes...)
	sign = append(sign, vBytes)
	return sign
}

//Verify Eth Signature
func VerifyEthSignature(ethtx *types.Transaction) error {
	sign := parseEthSignature(ethtx)
	if len(sign) <= 64 {
		return fmt.Errorf("eth signature lenght error:%v", len(sign))
	}

	signer := types.NewEIP155Signer(ethtx.ChainId())
	sighash := signer.Hash(ethtx)
	pub, err := crypto.Ecrecover(sighash[:], sign)
	if err != nil {
		return err
	}

	{

		signer := types.NewEIP155Signer(ethtx.ChainId())
		msg, _ := ethtx.AsMessage(signer, nil)

		p, err := crypto.SigToPub(sighash[:], sign)
		if err != nil {
			return err
		}

		addr := crypto.PubkeyToAddress(*p)
		if msg.From() != addr {
			return fmt.Errorf("verify sender failed! want:%v,got:%v", addr, msg.From())
		}
	}

	if !crypto.VerifySignature(pub, sighash[:], sign[:64]) {
		return fmt.Errorf("%v", "Verify Eth Signature failed")
	}
	return nil
}
