package token

import (
	"fmt"
	"testing"
)

func TestCreateJwtToen(t *testing.T) {
	tests := []struct {
		name       string
		wantJwtStr string
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "测试",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJwtStr, err := CreateJwtToken(1)
			fmt.Println(ParseJwtToken(gotJwtStr))
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateJwtToen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotJwtStr != tt.wantJwtStr {
				t.Errorf("CreateJwtToen() gotJwtStr = %v, want %v", gotJwtStr, tt.wantJwtStr)
			}
		})
	}
}
