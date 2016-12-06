# 有度企业应用Go sdk

## 功能介绍

- 企业应用http接口调用封装（sdk）
- 企业应用接口调用测试用例
- 企业应用回调接口测试用例
- AES_CBC加密测试用例

## 依赖库

- 无

## 开发工具

- LiteIDE 30.2
- Go 1.5或以上

## 加密算法

- 采用AES_CBC_256算法加密

## 测试方法

-切换到项目目录，运行go test -v命令，运行所有测试用例

## 集成方法
~~~ go
package main

import (
	"log"

	ydapp "github.com/youduim/EntAppSdkGo"
)

const (
	_Buin      = 666666
	_AppId     = `yd37D192E9F20E448192827A001A84D443`
	_EncAesKey = `AY2O92Lpyr4M2IOXT05NQG3eaXd72FlS/QZ1l4vGKsQ=`
	_User      = "sa08"
)

func main() {
	ydapp.Server_Addr := "http://localhost:7080" //设置服务器地址
	app, err := ydapp.NewMsgApp(_Buin, _AppId, _EncAesKey)
	if err != nil {
		log.Println("New app error:", err)
		return
	}
	token, expire, err := app.GetToken()
	if err != nil {
		log.Println("Get token error:", err)
		return
	}
	log.Printf("Token: %s, Expire: %d", token, expire) //expire为过期的时间戳，单位秒
	err = app.SendTxtMsg(_User, "123")
	if err != nil {
		log.Println("Send msg error:", err)
		return
	}
	log.Println("Send msg success.")
}
~~~



