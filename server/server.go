package main

import (
	. "EverySync/util"
	"github.com/xtaci/kcp-go"
	"log"
	"net"
)

func main(){
	server, err := kcp.ListenWithOptions(":10031", nil, 10, 3)
	CheckErr(err)
	for {
		client, err := server.Accept()
		CheckErr(err)
		go ReadDataFromClient(client)
	}
}


func ReadDataFromClient(client net.Conn){
	data := make([]byte, 1024)

	for{
		n, err := client.Read(data)
		CheckErr(err)
		log.Println("client msg:", string(data[:n]))
	}

}