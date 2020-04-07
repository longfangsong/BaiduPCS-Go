package main

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/handler"
	"github.com/iikira/BaiduPCS-Go/internal/pcsconfig"
	_ "github.com/iikira/BaiduPCS-Go/internal/pcsinit"
	"github.com/iikira/BaiduPCS-Go/pcsutil"
	"log"
	"net/http"
	"os"
)

func init() {
	pcsutil.ChWorkDir()

	err := pcsconfig.Config.Init()
	switch err {
	case nil:
	case pcsconfig.ErrConfigFileNoPermission, pcsconfig.ErrConfigContentsParseError:
		fmt.Fprintf(os.Stderr, "FATAL ERROR: config file error: %s\n", err)
		os.Exit(1)
	default:
		fmt.Printf("WARNING: config init error: %s\n", err)
	}
}

func main() {
	defer pcsconfig.Config.Close()
	http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	})
	http.HandleFunc("/file", handler.FileHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
