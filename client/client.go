package main

import (
	. "EverySync/util"
	"github.com/golang/protobuf/proto"
	"github.com/xtaci/kcp-go"
	"log"
	"net"
	"time"
	"EverySync/pb"
)

// const ServerIP = "149.28.59.218"
const ServerIP = "192.168.1.4"
const ServerPort  = 13301
func main(){

	dstAddr := &net.UDPAddr{IP: net.ParseIP(ServerIP), Port: ServerPort}
	srcAddr := &net.UDPAddr{IP: net.ParseIP(ServerIP), Port: 1223} // 注意端口必须固定
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	clientSess, err := kcp.NewConn(dstAddr.String(),nil,0,0,conn)

	//clientSess, err := kcp.DialWithOptions(dstAddr.String(),nil,0,0)

	CheckErr(err)
	log.Println("local addr : ", clientSess.LocalAddr().String())
	log.Println("Client connected to ", dstAddr.String())


	for {
		clientSess.Write([]byte("1111111111111111"))
		time.Sleep(1 * time.Second)
		log.Println("write success")
	}

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