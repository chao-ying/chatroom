package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"

	"go_code/chatroom/client/common/message"
)

//写一个函数，完成登陆
func login(userId int, userPwd string) (err error) {

	//下一个就要开始定协议
	// fmt.Printf("userId =%d userPwd=%s\n", userId, userPwd)
	// return nil

	//1.连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}

	//延时关闭
	defer conn.Close()

	//2.准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType //
	//3.创建一个LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	//4.将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//5.把data赋给mes.Data字段
	mes.data = string(data)

	//6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//7.到这个时候，data就是我们要发送的数据
	//7.1先把data的长度发送给服务器
	//先获取到
	// conn.Write(len(data))
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var bytes [4]byte
	binary.BigEndian.PutUint32(bytes[0:4], pkgLen)
	//发送长度
	n, err := conn.Write(bytes)
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	fmt.Println("客户端发送消息的长度ok")
}
