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

func main(){
	server, err := kcp.ListenWithOptions(":13301", nil, 10, 3)
	CheckErr(err)
	for {
		sess, err := server.Accept()
		CheckErr(err)
		go AddSessToCache(sess)
	}
}

// AddSessToCache
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

	var addrsList []string
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

	if len(addrsList) > 1 {
		log.Println("len : ", len(addrsList))
		go SendPeers(sess, id)
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