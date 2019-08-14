package api

// MsgBase msg base struct
type MsgBase struct {
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
}

// TextMsg text msg
type TextMsg struct {
	MsgBase
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	Safe          int `json:"safe"`
	EnableIDTrans int `json:"enable_id_trans"`
}

// ImageMsg image msg
type ImageMsg struct {
	MsgBase
	Image struct {
		MediaID string `json:"media_id"`
	} `json:"image"`
	Safe          int `json:"safe"`
	EnableIDTrans int `json:"enable_id_trans"`
}

// MarkdownMsg markdown msg
type MarkdownMsg struct {
	MsgBase
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

func NewTextMsg(content string) *TextMsg {
	msg := &TextMsg{MsgBase: MsgBase{MsgType: "text"}}
	msg.Text.Content = content
	return msg
}

func NewMarkdownMsg(content string) *MarkdownMsg {
	msg := &MarkdownMsg{MsgBase: MsgBase{MsgType: "markdown"}}
	msg.Markdown.Content = content
	return msg
}
