package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"

	"github.com/chao-ying/chatroom/common/message"
)

func readPkg(conn net.Conn) (mes message.Message, err error) {

	buf := make([]byte, 8096)
	fmt.Println("读取客户端发送的数据.......")
	//conn.Read在conn关闭的情况下，才会阻塞
	//如果客户端关闭了conn,则就不会阻塞
	_, err = conn.Read(buf[:4])
	//判断头消息
	if err != nil {
		fmt.Println("conn.Read  err=", err)
		//err = errors.New("read pkg header error")
		return
	}

	//根据buf[0:4]转成一个uint32的类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[:4])
	//根据pkgLen 读取消息内容
	n, err := conn.Read(buf[:pkgLen])
	//判读消息本身
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body error")
		fmt.Println("conn.Read fail err=", err)
		return
	}

	//把pkgLen反序列化 -》message.Message
	//技术就是一层窗户纸 &mes
	json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarsha err=", err)
		return
	}
	return
}

func writePkg(conn net.Conn, data []byte) (err error) {
	//	先把data的长度发送给对方

	var pkgLen uint32
	pkgLen = uint32(len(data))
	var bytes [4]byte //1 个int =4个byte

	binary.BigEndian.PutUint32(bytes[0:4], pkgLen)
	//发送长度
	n, err := conn.Write(bytes[0:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}

	//发送data本身
	n, err = conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	return

}
