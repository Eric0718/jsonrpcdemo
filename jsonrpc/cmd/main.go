package main

/*
	#ifdef __cplusplus
	extern "C" {
	#endif
	void* RunJsonRpc();
	void  StopJsonRpc(void *srv);
	#ifdef __cplusplus
	}
	#endif
*/
import "C"

import (
	"context"
	"jsonrpcdemo/jsonrpc/jsonrpcserver"
	"log"
	"net/http"
	"time"
	"unsafe"

	"github.com/mattn/go-pointer"
	"github.com/spf13/viper"
)

//export RunJsonRpc
func RunJsonRpc() unsafe.Pointer {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	viper := viper.New()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("ReadInConfig fail:", err.Error())
	}

	chainId := viper.GetString("chainId")
	netWorkId := viper.GetString("netWorkId")

	jsonrpcPort := viper.GetString("jsonrpcPort")
	archivePoint := viper.GetString("archivePoint")
	clinetVersion := viper.GetString("clinetVersion")

	srv := &http.Server{Addr: jsonrpcPort}
	s := jsonrpcserver.NewJsonRpcServer(chainId, netWorkId, archivePoint, clinetVersion)
	http.HandleFunc("/", s.HandRequest)

	go func() {
		log.Println("running jsonrpc server:", jsonrpcPort)
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("%v\n", err)
			return
		}
	}()
	return pointer.Save(srv)
}

//export StopJsonRpc
func StopJsonRpc(srv unsafe.Pointer) {
	s := pointer.Restore(srv).(*http.Server)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("StopJsonRpc failed:%v\n", err)
	}
}

func main() {}
