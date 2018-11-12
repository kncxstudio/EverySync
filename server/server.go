package main

import (
	"EverySync/pb"
	. "EverySync/util"
	"github.com/golang/protobuf/proto"
	"github.com/patrickmn/go-cache"
	"github.com/xtaci/kcp-go"
	"log"
	"net"
	"time"
)

var clientsCache = cache.New(5 * time.Minute, 10 * time.Minute)
var addrsCache = cache.New(5 * time.Minute, 10 * time.Minute)
var addrsList []string
const ListenPort = 13301

// var sessions = make([]*net.Conn,0)
func main(){

	laddr := &net.UDPAddr{IP:net.ParseIP("192.168.1.4"), Port: ListenPort}
	//raddr := &net.UDPAddr{IP:net.ParseIP("127.0.0.1"), Port: 9983}
	log.Println(laddr.String())
	conn, err := net.ListenUDP("udp", laddr)
	server, err := kcp.ServeConn(nil,0,0,conn)
	CheckErr(err)
	for {
		sess, err := server.Accept()
		CheckErr(err)
		log.Println(sess.RemoteAddr().String(), " connected!")
		go AddSessToCache(sess)
	}
}

// AddSessToCache add session to cache
func AddSessToCache(sess net.Conn){
	clientInfo := make([]byte,1024)
	n, _ := sess.Read(clientInfo)
	client := &pb.Client{}
	unmarshalErr := proto.Unmarshal(clientInfo[:n], client)
	CheckErr(unmarshalErr)
	client.Addr = sess.RemoteAddr().String()

	clientsCache.Set(sess.RemoteAddr().String(), client, cache.DefaultExpiration)
	id := client.GetID()

	temp, _ := addrsCache.Get(id)

	if temp == nil {
		addrsList = make([]string, 0)
	}else {
		addrsList = temp.([]string)
	}

	addrsList = append(addrsList, client.GetAddr())
	addrsCache.Set(id, addrsList, cache.DefaultExpiration)

	for _, addr := range addrsList{
		log.Println("addr : ", addr)
	}



	for {
		if len(addrsList) > 1 {
			log.Println("len : ", len(addrsList))
			go SendPeers(sess, id)
			break
		}
		time.Sleep(1 * time.Second)
	}

}



func SendPeers(sess net.Conn, id string){
	if addrs, found := addrsCache.Get(id); found {

		for _, addr := range addrs.([]string){
			if addr != sess.RemoteAddr().String(){
				client := &pb.Client{
					ID: id,
					Addr: addr,
				}
				data, err := proto.Marshal(client)
				CheckErr(err)
				sess.Write(data)
				break
			}
		}

	}

}