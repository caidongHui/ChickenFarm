package main

import (
	"bytes"
	"chickenFarm/db"
	"chickenFarm/model"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("信息处理错误")
		}
	}()
	db.Connct()
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	for {
		handleClient(conn)
	}
}
func handleClient(conn *net.UDPConn) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("数据处理错误")
		}
	}()
	var buf [512]byte
	msg, addr, err := conn.ReadFromUDP(buf[0:511])
	if err != nil {
		return
	}
	fmt.Println(msg, string(buf[0:]))
	daytime := "receive" //time.Now().String()
	var data model.UpInfo
	// 截取有效长度
	index := bytes.IndexByte(buf[0:], 0)
	if err := json.Unmarshal([]byte(string(buf[:index])), &data); err != nil {
		panic(err)
	}

	db.UpInsert(data)

	conn.WriteToUDP([]byte(daytime), addr)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
