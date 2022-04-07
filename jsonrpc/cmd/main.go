package main

import "C"

import (
	"context"
	"jsonrpcdemo/jsonrpc/jsonrpcserver"
	"jsonrpcdemo/jsonrpc/util"
	"jsonrpcdemo/logger"
	"net/http"
	"time"
	"unsafe"

	"github.com/mattn/go-pointer"
)

//export RunJsonRpc
func RunJsonRpc(chainId, netWorkId, archivePoint, clinetVersion, jsonrpcPort string) unsafe.Pointer {
	//通过make的方式，新构建一段内存来存放从C++处传入的字符串，深度拷贝防止C++中修改影响Go
	chid := util.MakeString(chainId)
	netid := util.MakeString(netWorkId)
	archivep := util.MakeString(archivePoint)
	cltv := util.MakeString(clinetVersion)
	jsonp := util.MakeString(jsonrpcPort)

	srv := &http.Server{Addr: jsonp}
	s := jsonrpcserver.NewJsonRpcServer(chid, netid, archivep, cltv)
	http.HandleFunc("/", s.HandRequest)

	go func() {
		logger.InitLogger()
		defer logger.SugarLogger.Sync()

		logger.SugarLogger.Infof("running jsonrpc server:%v", jsonp)
		err := srv.ListenAndServe()
		if err != nil {
			logger.SugarLogger.Warnf("%v", err)
			return
		}
	}()
	return pointer.Save(srv)
}

//export StopJsonRpc
func StopJsonRpc(srv unsafe.Pointer) C.int {
	s := pointer.Restore(srv).(*http.Server)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		logger.SugarLogger.Errorf("StopJsonRpc failed:%v", err)
		return C.int(-1)
	}
	return C.int(0)
}

func main() {}
