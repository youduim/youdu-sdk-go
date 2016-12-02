package ydapp

// 第三方接口URL
const (
	SCHEME = "http://"

	URL_YOUDU_API = "localhost:7080"

	API_GET_TOKEN = SCHEME + URL_YOUDU_API + "/v3/api/jgapp/ent.app.accesstoken.gen"

	API_SEND_MSG = SCHEME + URL_YOUDU_API + "/v3/api/jgapp/ent.app.msg.send"

	API_UPLOAD_FILE = SCHEME + URL_YOUDU_API + "/v3/api/jgapp/ent.app.media.upload"

	API_DOWNLOAD_FILE = SCHEME + URL_YOUDU_API + "/v3/api/jgapp/ent.app.media.get"

	API_SEARCH_FILE = SCHEME + URL_YOUDU_API + "/v3/api/jgapp/ent.app.media.search"

	CALLBACK_RECV_MSG = "/receive/youdu/msg"
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
