package main

import (
	"fmt"
	"io"
	"net"
)

//到的buf= [0 0 1 88]是字节数组  等于1*256+88 转换成int 类型是344
// func readPkg(conn net.Conn) (mes message.Message, err error) {

// 	buf := make([]byte, 8096)
// 	fmt.Println("读取客户端发送的数据.......")
// 	//conn.Read在conn关闭的情况下，才会阻塞
// 	//如果客户端关闭了conn,则就不会阻塞
// 	_, err = conn.Read(buf[:4])
// 	//判断头消息
// 	if err != nil {
// 		// fmt.Println("conn.Read  err=", err)
// 		//err = errors.New("read pkg header error")
// 		return
// 	}

// 	//根据buf[0:4]转成一个uint32的类型
// 	var pkgLen uint32
// 	pkgLen = binary.BigEndian.Uint32(buf[:4])
// 	//根据pkgLen 读取消息内容
// 	n, err := conn.Read(buf[:pkgLen])
// 	//判读消息本身
// 	if n != int(pkgLen) || err != nil {
// 		//err = errors.New("read pkg body error")
// 		fmt.Println("conn.Read fail err=", err)
// 		return
// 	}

// 	//把pkgLen反序列化 -》message.Message
// 	//技术就是一层窗户纸 &mes
// 	json.Unmarshal(buf[:pkgLen], &mes)
// 	if err != nil {
// 		fmt.Println("json.Unmarsha err=", err)
// 		return
// 	}
// 	return
// }

// func writePkg(conn net.Conn, data []byte) (err error) {
// 	//	先把data的长度发送给对方

// 	var pkgLen uint32
// 	pkgLen = uint32(len(data))
// 	var bytes [4]byte //1 个int =4个byte

// 	binary.BigEndian.PutUint32(bytes[0:4], pkgLen)
// 	//发送长度
// 	n, err := conn.Write(bytes[0:4])
// 	if n != 4 || err != nil {
// 		fmt.Println("conn.Write(bytes) fail", err)
// 		return
// 	}

// 	//发送data本身
// 	n, err = conn.Write(data)
// 	if n != int(pkgLen) || err != nil {
// 		fmt.Println("conn.Write(bytes) fail", err)
// 		return
// 	}
// 	return

// }

// //编写一个函数ServerProcessLogin函数，专门处理登陆请求
// func serverProcessLogin(conn net.Conn, mes *message.Message) (err error) {
// 	//核心代码
// 	//1.先从mes中取出mes.Data，并直接反序列化成LoginMes
// 	var loginMes message.LoginMes
// 	err = json.Unmarshal([]byte(mes.Data), &loginMes)
// 	if err != nil {
// 		fmt.Println("json.Unmarshal fail err=", err)
// 		return
// 	}
// 	//1.先声明一个resMes
// 	var resMes message.Message
// 	resMes.Type = message.LoginResMesType
// 	//2.再声明一个LoginResMes，并完成赋值
// 	var loginResMes message.LoginResMes

// 	//如果用户的id=100，密码=123456,认为合法，否则不合法
// 	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
// 		//合法
// 		loginResMes.Code = 200

// 	} else {
// 		//不合法
// 		loginResMes.Code = 500 //表示用户不存在
// 		loginResMes.Error = "该用户不存在,请注册再使用....."
// 	}

// 	//3。将loginResMes序列化
// 	data, err := json.Marshal(loginResMes)
// 	if err != nil {
// 		fmt.Println("json.Marshal fail err=", err)
// 		return
// 	}
// 	//4.将data赋值给resMes
// 	resMes.Data = string(data)

// 	//5.将resMes进行序列化，准备发送
// 	data, err = json.Marshal(resMes)
// 	if err != nil {
// 		fmt.Println("json.Marshal fail err=", err)
// 		return
// 	}
// 	//6.发送data,将其封装到writePkg函数中
// 	err = writePkg(conn, data)
// 	return
// }

//编写一个ServerProcessMes函数
//功能：根据客户端发送消息种类不同，决定调用哪个函数来处理
// func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {
// 	switch mes.Type {
// 	case message.LoginMesType:
// 		//处理登陆
// 		serverProcessLogin(conn, mes)
// 	case message.RegisterMesType:
// 		//处理注册
// 	default:
// 		fmt.Println("消息类型不存在，无法处理.....")
// 	}
// 	return
// }

//处理和 客户端的通讯
func process(conn net.Conn) {

	//这里需要及时关闭conn
	defer conn.Close()

	//循环的读客户端发送的信息
	for {

		//这里我们将读取数据包，直接封装成一个函数readPkg(),返回Message,Err
		mes, err := readPkg(conn)
		//读取客户端发送的消息

		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器端也退出......")
				return
			}
			fmt.Println("readPkg err=", err)

		}
		// fmt.Println("mes=", mes)
		err = serverProcessMes(conn, &mes)
		if err != nil {
			return
		}

		// fmt.Println("读到的buf=", buf[:4])
	}

}
func main() {
	// 提示信息
	fmt.Println("服务器在8889端口监听......")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	//一旦监听成功，就等待客户端来链接服务器
	for {
		fmt.Println("等待客户端来连接服务器")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}

		//一旦链接成功，则启动一个协程和客户端保持通讯
		go process(conn) // 它里面的defer conn.close()并没有关闭，只有在执行新的代码时候就会来关闭
		fmt.Println("hello world")

	}
}
