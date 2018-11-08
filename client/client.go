package main

import (
	"EverySync/pb"
	. "EverySync/util"
	"github.com/golang/protobuf/proto"
	"github.com/xtaci/kcp-go"
	"log"
	"net"
	"time"
)

// const ServerIP = "149.28.59.218"
const ServerIP = "127.0.0.1"
const ServerPort  = 13301
func main(){

	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 9983} // 注意端口必须固定
	dstAddr := &net.UDPAddr{IP: net.ParseIP(ServerIP), Port: ServerPort}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	clientSess, err := kcp.NewConn(dstAddr.String(),nil,10,3,conn)

	CheckErr(err)
	log.Println("Client connected to ", dstAddr.String())


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


	sess, err := kcp.DialWithOptions(otherClient.Addr, nil, 10, 3)
	sess.Write([]byte("dig msg"))
	go ReadFromSession(sess)

	for {
		sess.Write([]byte("data msg"))
		time.Sleep(2 * time.Second)
	}
}


func ReadFromSession(session net.Conn){
	data := make([]byte,1024)

	for {
		n, readErr := session.Read(data)
		CheckErr(readErr)
		log.Println(string(data[:n]))
	}
}