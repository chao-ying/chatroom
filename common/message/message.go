package message

const (
	LoginMesType    = "LoginMes"
	LoginResMesType = "LoginResMes"
)

type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息的类型
}

//定义两个消息，后面需要在增加
type LoginMes struct {
	UserId   int    `json:"userId"`   //用户密码
	UserPwd  string `json:"userPwd"`  //用户ID
	UserName string `json:"userName"` //用户名

}

type LoginResMes struct {
	Code  int    `json:"code"`  //返回状态码500表示该用户	注册，200表示登陆成功
	Error string `json:"error"` //返回错误消息
}
