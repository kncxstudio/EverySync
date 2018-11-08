package main

import (
	"EverySync/pb"
	. "EverySync/util"
	"github.com/golang/protobuf/proto"
	"github.com/xtaci/kcp-go"
	"log"
)

const ServerAddr = "149.28.59.218:13301"
func main(){
	clientSess, err := kcp.DialWithOptions(ServerAddr, nil, 10, 3)
	CheckErr(err)
	log.Println("Client connected to ", ServerAddr)


	client := &pb.Client{
		ID: "asdadawfgd",
		Addr: "",
	}

	data, _ := proto.Marshal(client)
	clientSess.Write(data)


	data = make([]byte, 1024)
	n, err := clientSess.Read(data)
	CheckErr(err)
	otherClient := &pb.Client{}
	unmarshalErr := proto.Unmarshal(data[:n], otherClient)
	CheckErr(unmarshalErr)

	log.Println(otherClient)
	//for {
	//	_, err := clientSess.Write([]byte("clientSess says hello!"))
	//	CheckErr(err)
	//	log.Println("client write success!")
	//	time.Sleep(1 * time.Second)
	//}
}
