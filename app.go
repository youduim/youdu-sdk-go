package ydapp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	HttpJsonType = "application/json"
)

var (
	Buin      int32
	AppId     string
	EncAesKey string
)

type Receiver interface {
	Receive(*ReceiveMsg)
}

type MsgApp struct {
	buin     int32
	aesKey   []byte
	appId    string
	accToken string
	recv     Receiver
	hc       *http.Client
}

/*
	@buin 企业号
	@appId 应用的id
	@encAesKey base64编码后的AesKey(256位长度)
*/
func NewMsgApp(buin int32, appId, encAesKey string) (*MsgApp, error) {
	key, err := base64.StdEncoding.DecodeString(encAesKey)
	if err != nil {
		return nil, errors.New("Base64 decode error: " + err.Error())
	}
	if len(key) != 32 {
		return nil, errors.New("Invalid aes key size")
	}
	return &MsgApp{
		buin:   buin,
		aesKey: key,
		appId:  appId,
		hc:     &http.Client{},
	}, nil
}

/*
	设置回调
	需监听一个回调的端口
	如：http.ListenAndServe(":8899", demo)

*/
func (this *MsgApp) SetReceiver(r Receiver) {
	this.recv = r
}

func (this *MsgApp) encrypt(data []byte) (string, error) {
	return AesEncrypt(data, this.aesKey, this.appId)
}

func (this *MsgApp) decrypt(s string) (*RawMsg, error) {
	return AesDecrypt(s, this.aesKey)
}

func (this *MsgApp) post(api, ct string, req []byte) (*ApiResponse, error) {
	httpRsp, err := this.hc.Post(Server_Addr+api+"?accessToken="+this.accToken, ct, bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}
	if httpRsp.StatusCode != http.StatusOK {
		return nil, Error(httpRsp.Status, nil)
	}
	bd, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}

	rsp, err := NewReponse(bd)
	if err != nil {
		return nil, Error("Response unmarshal error", err)
	}
	return rsp, nil
}

func (this *MsgApp) getFile(req []byte) ([]byte, error) {
	httpRsp, err := this.hc.Post(Server_Addr+API_DOWNLOAD_FILE+"?accessToken="+this.accToken, HttpJsonType, bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}
	if httpRsp.StatusCode != http.StatusOK {
		return nil, Error(httpRsp.Status, nil)
	}
	bd, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}
	return bd, nil
}

/*
	获取token
	加密一个时间戳，传到服务器
*/
func (this *MsgApp) GetToken() (string, int64, error) {
	timex := fmt.Sprint(time.Now().Unix())
	cipherText, err := this.encrypt([]byte(timex))
	if err != nil {
		return "", 0, Error("Encrypt error", err)
	}

	req := NewRequest()
	req.Set("buin", this.buin)
	req.Set("appId", this.appId)
	req.Set("encrypt", cipherText)
	bs, _ := req.Encode()

	rsp, err := this.post(API_GET_TOKEN, HttpJsonType, bs)
	if err != nil {
		return "", 0, Error("Post to get token error", err)
	}
	if !rsp.StatusOK() {
		return "", 0, rsp.Error()
	}

	enc, err := rsp.GetString("encrypt")
	if err != nil {
		return "", 0, Error("Get body error", err)
	}
	raw, err := this.decrypt(enc)
	if err != nil {
		return "", 0, Error("Decrypt access token error:", err)
	}
	js, _ := NewJson(raw.Data)
	this.accToken, _ = js.Get("accessToken").String()

	expire, _ := js.Get("expireIn").Int64()
	return this.accToken, expire, nil
}

/*
	上传图片
	传入图片名字与图片数据
	支持jpg, png, gif格式
*/
func (this *MsgApp) UploadImageBytes(name string, data []byte) (string, error) {
	return this.upload(MediaTypeImage, name, data)
}

/*
	上传图片
	传入图片名字与路径
	支持jpg, png, gif格式
*/
func (this *MsgApp) UploadImage(name string, path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return this.upload(MediaTypeImage, name, data)
}

/*
	上传文件
	传入文件名与文件数据
*/
func (this *MsgApp) UploadFileBytes(name string, data []byte) (string, error) {
	return this.upload(MediaTypeFile, name, data)
}

/*
	上传文件
	传入文件名与文件路径
*/
func (this *MsgApp) UploadFile(name string, path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return this.upload(MediaTypeFile, name, data)
}

func (this *MsgApp) upload(ftype, fname string, data []byte) (string, error) {
	body := bytes.NewBufferString("")
	mwr := multipart.NewWriter(body)
	req := NewRequest()
	req.Set("type", ftype)
	req.Set("name", fname)

	bs, _ := req.Encode()
	enc, _ := this.encrypt(bs)
	mwr.WriteField("encrypt", enc)

	pw, _ := mwr.CreateFormFile("file", fname)
	msg, _ := this.encrypt(data)
	pw.Write([]byte(msg))
	mwr.Close()

	rsp, err := this.post(API_UPLOAD_FILE, mwr.FormDataContentType(), body.Bytes())
	if err != nil {
		return "", Error("Post to upload file error", err)
	}
	if !rsp.StatusOK() {
		return "", rsp.Error()
	}

	enc, err = rsp.GetString("encrypt")
	if err != nil {
		return "", Error("Get encrypt error", err)
	}
	raw, err := this.decrypt(enc)
	if err != nil {
		return "", Error("Aes decrypt error", err)
	}
	pm, err := NewJson(raw.Data)
	if err != nil {
		return "", Error("Json unmarshal error", err)
	}
	mediaId, _ := pm.Get("mediaId").String()
	return mediaId, nil
}

/*
	下载文件
	传入mediaId
	返回文件数据
*/
func (this *MsgApp) DownloadFile(mediaId string) ([]byte, error) {
	bs, err := this.download(mediaId)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

/*
	下载文件
	并保存到指定路径
	自动创建路径中的目录与文件
*/
func (this *MsgApp) DownloadFileSave(mediaId string, path string) error {
	data, err := this.download(mediaId)
	if err != nil {
		return err
	}
	return this.save(data, path)
}

/*
	下载图片
	返回图片数据
*/
func (this *MsgApp) DownloadImage(mediaId string) ([]byte, error) {
	bs, err := this.download(mediaId)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

/*
	下载图片
	并保存到指定路径
	自动创建路径中的目录与文件
*/
func (this *MsgApp) DownloadImageSave(mediaId string, path string) error {
	data, err := this.download(mediaId)
	if err != nil {
		return err
	}
	return this.save(data, path)
}

func (*MsgApp) save(data []byte, path string) error {
	dir, _ := filepath.Split(path)
	err := os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0666)
}

func (this *MsgApp) download(mediaId string) ([]byte, error) {
	req := NewRequest()
	req.Set("mediaId", mediaId)

	bs, _ := req.Encode()
	enc, _ := this.encrypt(bs)
	em := NewRequest()
	em.Set("encrypt", enc)

	bs, _ = em.Encode()
	bd, err := this.getFile(bs)
	if err != nil {
		return nil, Error("Download file error", err)
	}

	raw, err := this.decrypt(string(bd))
	if err != nil {
		return nil, Error("Decrypt error", err)
	}
	return raw.Data, nil
}

func (this *MsgApp) SearchFile(mediaId string) (string, int64, error) {
	req := NewRequest()
	req.Set("mediaId", mediaId)

	bs, _ := req.Encode()
	enc, _ := this.encrypt(bs)
	em := NewRequest()
	em.Set("encrypt", enc)

	bs, _ = em.Encode()
	rsp, err := this.post(API_SEARCH_FILE, HttpJsonType, bs)
	if err != nil {
		return "", 0, Error("Post to search file error", err)
	}
	enc, err = rsp.GetString("encrypt")
	if err != nil {
		return "", 0, Error("Get encrypt error", err)
	}
	raw, err := this.decrypt(enc)
	if err != nil {
		return "", 0, Error("Decrypt error", err)
	}
	pm, _ := NewJson(raw.Data)
	name, err := pm.Get("name").String()
	if err != nil {
		return "", 0, err
	}
	size, err := pm.Get("size").Int64()
	if err != nil {
		return "", 0, err
	}
	return name, size, nil
}

/*
	发送文本消息(包括表情)
	传入接收者用户名与消息内容
	如果发送给多人，则用户间用"|"隔开，如cs1|cs2|cs3
*/
func (this *MsgApp) SendTxtMsg(toUser, content string) error {
	msg := NewRequest()
	msg.Set("toUser", toUser)
	msg.Set("msgType", MsgTypeText)
	msg.Set("text", map[string]interface{}{
		"content": content,
	})

	bs, _ := msg.Encode()
	enc, _ := this.encrypt(bs)
	req := NewRequest()
	req.Set("buin", this.buin)
	req.Set("appId", this.appId)
	req.Set("encrypt", enc)

	bs, _ = req.Encode()
	rsp, err := this.post(API_SEND_MSG, HttpJsonType, bs)
	if err != nil {
		return Error("Send text msg error", err)
	}
	if !rsp.StatusOK() {
		return rsp.Error()
	}
	return nil
}

/*
	发送图片消息
	传入已上传图片的mediaId与接收者用户名
	如果发送给多人，则用户间用"|"隔开，如cs1|cs2|cs3
*/

func (this *MsgApp) SendImg(toUser, path string) error {
	_, name := filepath.Split(path)
	mediaId, err := this.UploadFile(name, path)
	if err != nil {
		return err
	}
	return this.SendImgMsg(toUser, mediaId)
}

func (this *MsgApp) SendImgMsg(toUser, mediaId string) error {
	msg := NewRequest()
	msg.Set("toUser", toUser)
	msg.Set("msgType", MsgTypeImage)
	msg.Set("image", map[string]interface{}{
		"media_id": mediaId,
	})
	bs, _ := msg.Encode()
	enc, _ := this.encrypt(bs)

	req := NewRequest()
	req.Set("buin", this.buin)
	req.Set("appId", this.appId)
	req.Set("encrypt", enc)

	bs, _ = req.Encode()
	rsp, err := this.post(API_SEND_MSG, HttpJsonType, bs)
	if err != nil {
		return Error("Send image msg error", err)
	}
	return rsp.Error()
}

/*
	发送文件消息
	传入已上传文件的mediaId与接收者用户名
	如果发送给多人，则用户间用"|"隔开，如cs1|cs2|cs3
*/
func (this *MsgApp) SendFile(toUser, name, path string) error {
	mediaId, err := this.UploadFile(name, path)
	if err != nil {
		return err
	}
	return this.SendFileMsg(toUser, mediaId)
}

func (this *MsgApp) SendFileMsg(toUser, mediaId string) error {
	msg := NewRequest()
	msg.Set("toUser", toUser)
	msg.Set("msgType", MsgTypeFile)
	msg.Set("file", map[string]interface{}{
		"media_id": mediaId,
	})

	bs, _ := msg.Encode()
	enc, _ := this.encrypt(bs)
	req := NewRequest()
	req.Set("buin", this.buin)
	req.Set("appId", this.appId)
	req.Set("encrypt", enc)

	bs, _ = req.Encode()
	rsp, err := this.post(API_SEND_MSG, HttpJsonType, bs)
	if err != nil {
		return Error("send file msg error", err)
	}
	return rsp.Error()
}

/*
	发送图片文章消息
	@title 文章标题
	@media_id 图片的media_id, 如果media_id为空, 则从path读取文件
	@digest 文章摘要
	@content 文章正文内容
	@showFront 是否在正文显示图片
	@toUser 消息接收者
	如果发送给多人，则用户间用"|"隔开，如cs1|cs2|cs3
*/
func (this *MsgApp) SendMpnewsMsg(toUser string, mpnews []*Mpnews) error {
	mplist := make([]interface{}, 0)
	for _, mp := range mpnews {
		news := NewRequest()
		news.Set("title", mp.Title)
		if len(mp.MediaId) == 0 {
			mp.MediaId, _ = this.UploadImage("mpnews.jpg", mp.Path)
		}
		news.Set("media_id", mp.MediaId)
		news.Set("digest", mp.Digest)
		news.Set("content", mp.Content)
		news.Set("url", mp.Url)
		news.Set("showFront", mp.ShowFront)
		mplist = append(mplist, news)
	}
	msg := NewRequest()
	msg.Set("toUser", toUser)
	msg.Set("msgType", MsgTypeMpnews)
	msg.Set("mpnews", mplist)

	bs, _ := msg.Encode()
	enc, _ := this.encrypt(bs)
	req := NewRequest()
	req.Set("buin", this.buin)
	req.Set("appId", this.appId)
	req.Set("encrypt", enc)
	bs, _ = req.Encode()
	rsp, err := this.post(API_SEND_MSG, HttpJsonType, bs)
	if err != nil {
		return Error("Send mpnews error", err)
	}
	return rsp.Error()
}

/*
	发送外链信息
	@title 标题
	@url 要跳转到的链接
	@digest 文章摘要
	@media_id 图片的media_id, 如果media_id为空, 则从path读取文件
*/
func (this *MsgApp) SendExlinkMsg(toUser string, links []*Exlink) error {
	list := make([]interface{}, 0)
	for _, link := range links {
		ex := NewRequest()
		ex.Set("title", link.Title)
		ex.Set("url", link.Url)
		ex.Set("digest", link.Digest)
		if len(link.MediaId) == 0 {
			link.MediaId, _ = this.UploadImage("exlink.jpg", link.Path)
		}
		ex.Set("media_id", link.MediaId)
		list = append(list, ex)
	}
	msg := NewRequest()
	msg.Set("toUser", toUser)
	msg.Set("msgType", MsgTypeExlink)
	msg.Set("exlink", list)

	bs, _ := msg.Encode()
	enc, _ := this.encrypt(bs)
	req := NewRequest()
	req.Set("buin", this.buin)
	req.Set("appId", this.appId)
	req.Set("encrypt", enc)

	bs, _ = req.Encode()
	rsp, err := this.post(API_SEND_MSG, HttpJsonType, bs)
	if err != nil {
		return err
	}
	return rsp.Error()
}

func (this *MsgApp) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	this.Receive(rw, req)
}

/*
	回调消息
	接收回调需设置Receiver
	并且需设置监听端口，例
	http.ListenAndServe(":8899", app)
*/
func (this *MsgApp) Receive(rw http.ResponseWriter, req *http.Request) {
	log.Println("Receive msg")
	bs, _ := ioutil.ReadAll(req.Body)
	var p ReceivePack
	err := json.Unmarshal(bs, &p)
	if err != nil {
		log.Println("Receive package error:", err)
		return
	}
	raw, err := this.decrypt(p.Encrypt)
	if err != nil {
		log.Println("Decrypt error:", err)
		return
	}

	var msg ReceiveMsg
	err = json.Unmarshal(raw.Data, &msg)
	if err != nil {
		log.Println("Json unmarshal error:", err)
		return
	}

	if this.recv != nil {
		this.recv.Receive(&msg)
	}

	log.Println("Recv msg:", msg)
	rw.Write([]byte(msg.PackageId))
}
