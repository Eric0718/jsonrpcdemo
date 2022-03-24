package main

import "C"

import (
	"jsonrpcdemo/jsonRPC"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func main() {}

//export RunServer
func RunServer() {
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

	listenPort := viper.GetString("listenPort")
	consensusPoint := viper.GetString("consensusPoint")
	archivePoint := viper.GetString("archivePoint")
	clinetVersion := viper.GetString("clinetVersion")

	s := jsonRPC.NewServer(chainId, netWorkId, consensusPoint, archivePoint, clinetVersion)
	http.HandleFunc("/", s.HandRequest)

	log.Println("running jsonrpc server:", listenPort)
	err := http.ListenAndServe(listenPort, nil)
	if err != nil {
		log.Fatalf("Start server failed:%v\n", err)
	}
}
