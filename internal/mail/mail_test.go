package mail

import (
	"testing"
)

func TestSend(t *testing.T) {
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
			err := Send(tt.email, tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
