package ydapp

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
)

const (
	_Buin      = 666666
	_AppId     = `yd37D192E9F20E448192827A001A84D443`
	_EncAesKey = `AY2O92Lpyr4M2IOXT05NQG3eaXd72FlS/QZ1l4vGKsQ=`
	_User      = "sa08"
)

func TestAllApi(t *testing.T) {
	demo, _ := NewMsgApp(_Buin, _AppId, _EncAesKey)

	go http.ListenAndServe(":8899", demo)

	//获取token
	_, _, err := demo.GetToken()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Get accesstoken success:", demo.accToken)

	//传入路径，上传文件
	fileId, err := demo.UploadFile("hello.txt", filepath.Join("file", "hello.txt"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("File mediaId:", fileId)

	//传入路径，上传图片
	imgId, err := demo.UploadImage("lake.jpg", filepath.Join("file", "lake.jpg"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Image mediaId:", imgId)

	//下载图片并保存到指定路径
	err = demo.DownloadFileSave(fileId, "file/file.txt")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Download and save file success.")

	//下载文件，返回文件数据
	data, err := demo.DownloadFile(fileId)
	if err != nil {
		t.Error(err)
		return
	}
	ioutil.WriteFile("file/file1.txt", data, 0666)
	t.Log("Download file success.")

	//下载文件并保存到指定路径
	err = demo.DownloadImageSave(imgId, "file/image.jpg")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Download and save image success.")

	//下载图片，返回图片数据
	data, err = demo.DownloadImage(imgId)
	if err != nil {
		t.Error(err)
	}
	ioutil.WriteFile("file/image1.jpg", data, 0666)
	t.Log("Download file success.")

	//发送文本消息
	err = demo.SendTxtMsg(_User, "第三方接口测试123abc!@#$%^^&*()/::|/::)-+=-")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Send text msg success.")

	//传入mediaId，发送图片信息
	err = demo.SendImgMsg(_User, imgId)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Send image msg success.")

	//传入mediaId，发送文件信息
	err = demo.SendFileMsg(_User, fileId)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Send file msg success.")

	//传入路径，上传并发送图片信息
	err = demo.SendImg(_User, "file/lake.jpg")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Send image success.")

	//传入路径，上传并发送文件信息
	err = demo.SendFile(_User, "hello.txt", "file/hello.txt")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Send file success.")

	m := &Mpnews{
		Title: "测试标题1",
		Path:  "file/lake.jpg",
		//MediaId:   imgId,
		Digest:    "一些摘要",
		Url:       `http://www.zhbuswx.com/busline/BusQuery.html?v=1.97#/`,
		Content:   "我本将心向明月，奈何明月照沟渠",
		ShowFront: 1,
	}
	//发送图文信息
	err = demo.SendMpnewsMsg(_User, []*Mpnews{m})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Send msg success.")

	link := &Exlink{
		Title:  "这是一个标题",
		Url:    "http://www.zhbuswx.com/busline/BusQuery.html?v=1.97#/",
		Digest: "外链摘要",
		Path:   "file/lake.jpg",
		//MediaId: imgId,
	}
	err = demo.SendExlinkMsg(_User, []*Exlink{link})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Send exlink msg success.")

	exist, err := demo.SearchFile(fileId)
	if err != nil {
		t.Error(err)
	}
	if !exist {
		t.Error("File not exist")
	}
	t.Log("Search file success.")
}
