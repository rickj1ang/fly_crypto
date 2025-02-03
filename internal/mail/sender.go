package mail

type Message struct {
	SendTo      string
	TargetPrice float64
	CoinSymbol  string
}

func Sender(messages chan Message) {
	for message := range messages {
		SendNotify(message)
	}
}
