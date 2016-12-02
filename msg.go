package ydapp

type Mpnews struct {
	Title     string
	MediaId   string
	Path      string
	Digest    string
	Url       string
	Content   string
	ShowFront int32
}

type Exlink struct {
	Title   string
	Url     string
	Digest  string
	MediaId string
	Path    string
}

type ReceivePack struct {
	ToBuin  int32  `json:"toBuin"`
	ToApp   string `json:"toApp"`
	Encrypt string `json:"encrypt"`
}

type ReceiveMsg struct {
	Buin       int32                  `json:"buin"`
	FromUser   string                 `json:"fromUser"`
	FromGid    int64                  `json:"fromGid"`
	CreateTime int64                  `json:"createTime"`
	PackageId  string                 `json:"packageId"`
	MsgType    string                 `json:"msgType"`
	SessionId  string                 `json:"sessionId"`
	Text       map[string]interface{} `json:"text"`
	Image      map[string]interface{} `json:"image"`
	File       map[string]interface{} `json:"file"`
}
