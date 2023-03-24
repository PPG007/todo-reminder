package util

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"path/filepath"
	"todo-reminder/proto"
	"todo-reminder/repository/bsoncodec"
)

var (
	client proto.ChatGPTServiceClient
)

func init() {
	if !viper.GetBool("chatgpt.enabled") {
		return
	}
	conn, err := grpc.Dial(viper.GetString("chatgpt.url"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client = proto.NewChatGPTServiceClient(conn)
}

func GetImageFile(ctx context.Context, text string) (string, string, error) {
	resp, err := client.GetImageResponse(ctx, &proto.String{
		Value: text,
	})
	if err != nil {
		return "", "", err
	}
	name := fmt.Sprintf("%s.jpg", bsoncodec.NewObjectId().Hex())
	file, err := os.Create(name)
	if err != nil {
		return "", "", err
	}
	_, err = file.Write(resp.Data)
	if err != nil {
		return "", "", err
	}
	file.Close()
	abs, err := filepath.Abs(file.Name())
	if err != nil {
		return "", "", err
	}
	return abs, name, nil
}

func GetTextResponse(ctx context.Context, input string) (string, error) {
	resp, err := client.GetTextResponse(ctx, &proto.String{
		Value: input,
	})
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}
