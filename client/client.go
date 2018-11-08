package main

import (
	."EverySync/util"
	"github.com/xtaci/kcp-go"
	"log"
	"time"
)

const ServerAddr = "localhost:10031"
func main(){
	clientSess, err := kcp.DialWithOptions(ServerAddr, nil, 10, 3)
	CheckErr(err)
	log.Println("Client connected to ", ServerAddr)

	for {
		_, err := clientSess.Write([]byte("clientSess says hello!"))
		CheckErr(err)
		log.Println("client write success!")
		time.Sleep(1 * time.Second)
	}
}
