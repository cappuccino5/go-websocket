package chat

const (
	// 文本消息
	Type_Text = iota

	// 图片消息
	Type_Image

	// 音频消息
	Type_Audio

	// 首条消息
	Type_First

	//用户注册消息[提交用户的user id]
	Type_Register
)

type Message struct {
	From     uint64    `json:"from"`
	To       uint64    `json:"to"`
	Type     uint8     `json:"type"`
	Time     uint64    `json:"create_time"`
	MsgIndex uint64    `json:"msg_index"`
	Content  []Content `json:"content"`
}

type Content struct {
	Text     string `json:"text"` // 内容
	Url      string `json:"url"`  // 消息的url,文字消息的话则无
	ThumbUrl string `json:"thumb_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
	Duration int    `json:"duration"`
}
