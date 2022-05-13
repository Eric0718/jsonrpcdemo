package main

import (
	"jsonrpcdemo/jsonrpc/jsonrpcserver"
	"jsonrpcdemo/logger"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

/*
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
*/
func main() {
	viper := viper.New()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("ReadInConfig fail:", err.Error())
		os.Exit(1)
	}

	chainId := viper.GetString("chainId")
	netWorkId := viper.GetString("netWorkId")

	jsonrpcPort := viper.GetString("jsonrpcPort")
	//consensusPoint := viper.GetString("consensusPoint")
	archivePoint := viper.GetString("archivePoint")
	clinetVersion := viper.GetString("clinetVersion")

	{
		logger.InitLogger()
		defer logger.SugarLogger.Sync()

		logger.SugarLogger.Infof("running jsonrpc server:%v,chainid:%v", jsonrpcPort, chainId)
	}

	srv := &http.Server{Addr: jsonrpcPort}
	s := jsonrpcserver.NewJsonRpcServer(chainId, netWorkId, archivePoint, clinetVersion)
	http.HandleFunc("/", s.HandRequest)
	err := srv.ListenAndServe()
	if err != nil {
		logger.SugarLogger.Warnf("%v", err)
		return
	}

}
