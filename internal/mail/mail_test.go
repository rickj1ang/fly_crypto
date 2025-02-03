package mail

import (
	"testing"
	"time"
)

func TestSendVerifyCode(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		code    string
		wantErr bool
	}{
		{
			name:    "valid email test",
			email:   "jdr666@foxmail.com",
			code:    "123456",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendVerifyCode(tt.email, tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendNotify(t *testing.T) {
	notifys := []Message{
		{
			SendTo:      "jdr666@foxmail.com",
			TargetPrice: 75000.5,
			CoinSymbol:  "BTC",
		},
		{
			SendTo:      "jdr666@foxmail.com",
			TargetPrice: 85000.5,
			CoinSymbol:  "BTC",
		},
		{
			SendTo:      "jdr666@foxmail.com",
			TargetPrice: 95000.5,
			CoinSymbol:  "BTC",
		},
	}

	messages := make(chan Message, 10)
	go Sender(messages)
	for _, notify := range notifys {
		messages <- notify
	}
	time.Sleep(30 * time.Second)
}
