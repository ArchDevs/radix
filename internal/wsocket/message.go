package wsocket

type IncomingMessage struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type OutgoingMessage struct {
	From      string `json:"from"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}
