package util

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"github.com/spf13/viper"
	"time"
)

var (
	encryptor cipher.Block
	days      int
)

type Token struct {
	UserId    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

func init() {
	c, err := aes.NewCipher([]byte(viper.GetString("token.key")))
	if err != nil {
		panic(err)
	}
	encryptor = c
	days = viper.GetInt("token.validDays")
	if days == 0 {
		panic("need token valid days")
	}
}

func GenerateToken(userId string) string {
	token := Token{
		UserId:    userId,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().AddDate(0, 0, days),
	}
	source := []byte(MarshalToJson(token))
	out := make([]byte, 0, len(source))
	encryptor.Encrypt(out, source)
	return string(out)
}

func ParseToken(ctx context.Context, tokenStr string) (*Token, error) {
	source := []byte(tokenStr)
	out := make([]byte, 0, len(source))
	encryptor.Decrypt(out, source)
	token := MustUnmarshalFromJson[Token](string(out))
	if time.Now().After(token.ExpiredAt) {
		return nil, errors.New("invalid token")
	}
	return &token, nil
}
