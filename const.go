package ydapp

var Server_Addr string                         //服务器地址和端口(协议务必带上)，例: http://localhost:7080
var Callback_Url string = "/receive/youdu/msg" //设置回调的URI

// 第三方接口URL
var (
	API_GET_TOKEN = "/cgi/gettoken"

	API_SEND_MSG = "/cgi/msg/send"

	API_UPLOAD_FILE = "/cgi/media/upload"

	API_DOWNLOAD_FILE = "/cgi/media/get"

	API_SEARCH_FILE = "/cgi/media/search"
)

//文件类型定义
const (
	MediaTypeFile  = "file"  //文件
	MediaTypeImage = "image" //图片
)

//消息类型定义
const (
	MsgTypeText   = "text"   //文本
	MsgTypeFile   = "file"   //文件
	MsgTypeImage  = "image"  //图片
	MsgTypeMpnews = "mpnews" //图文
	MsgTypeExlink = "exlink" //外链
)
