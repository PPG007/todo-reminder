package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
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
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)]
}

func Encrypt(data []byte) []byte {
	data = pkcs7Padding(data, encryptor.BlockSize())
	encrypted := make([]byte, len(data))
	cipher.NewCBCEncrypter(encryptor, []byte(viper.GetString("token.key"))[:encryptor.BlockSize()]).CryptBlocks(encrypted, data)
	return encrypted
}

func Decrypt(data []byte) []byte {
	decrypted := make([]byte, len(data))
	cipher.NewCBCDecrypter(encryptor, []byte(viper.GetString("token.key"))[:encryptor.BlockSize()]).CryptBlocks(decrypted, data)
	decrypted = pkcs7UnPadding(decrypted)
	return decrypted
}

func GenerateToken(userId string) string {
	token := Token{
		UserId:    userId,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().AddDate(0, 0, days),
	}
	return base64.StdEncoding.EncodeToString(Encrypt([]byte(MarshalToJson(token))))
}

func ParseToken(tokenStr string) (*Token, error) {
	data, err := base64.StdEncoding.DecodeString(tokenStr)
	if err != nil {
		return nil, err
	}
	token := MustUnmarshalFromJson[Token](string(Decrypt(data)))
	if token.CreatedAt.Unix() > 0 && token.ExpiredAt.After(time.Now()) {
		return &token, nil
	}
	return nil, errors.New("invalid token")
}
