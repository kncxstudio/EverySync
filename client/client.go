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
const ServerIP = "127.0.0.1"
const ServerPort  = 13301
const LocalPort = 1223
func main(){

	dstAddr := &net.UDPAddr{IP: net.ParseIP(ServerIP), Port: ServerPort}
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: LocalPort} // 注意端口必须固定
	//conn, err := net.DialUDP("udp", srcAddr, dstAddr)

	//conn, err := kcp.DialWithOptions(dstAddr.String(),nil,0,0)
	udpConn, err := net.ListenUDP("udp", srcAddr)
	conn, err := kcp.NewConn(dstAddr.String(),nil,0,0,udpConn)
	CheckErr(err)
	log.Println("Client connected to ", dstAddr.String())

	client := &pb.Client{
		ID: "asdadawfgd",
		Addr: "",
	}
	data, _ := proto.Marshal(client)
	conn.Write(data)


	data = make([]byte, 1024)
	n, err := conn.Read(data)
	CheckErr(err)
	conn.Close()



	otherClient := &pb.Client{}
	unmarshalErr := proto.Unmarshal(data[:n], otherClient)
	CheckErr(unmarshalErr)
	log.Println(otherClient)
	rAddr, err := net.ResolveUDPAddr("udp", otherClient.Addr)
	CheckErr(err)
	udpConn, listenErr := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: LocalPort})
	CheckErr(listenErr)
	sess, err := kcp.NewConn(rAddr.String(), nil, 0, 0,udpConn)
	sess.Write([]byte("dig msg"))
	go ReadFromSession(sess)

	for {
		sess.Write([]byte("data msg"))
		time.Sleep(2 * time.Second)
		log.Println(sess.LocalAddr(), " writes msg")
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