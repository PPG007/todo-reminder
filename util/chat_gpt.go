package util

import (
	"context"
	"fmt"
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
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client = proto.NewChatGPTServiceClient(conn)
}

func GetImageFile(ctx context.Context, text string) (string, error) {
	resp, err := client.GetImageResponse(ctx, &proto.String{
		Value: text,
	})
	if err != nil {
		return "", err
	}
	file, err := os.Create(fmt.Sprintf("%s.jpg", bsoncodec.NewObjectId().Hex()))
	if err != nil {
		return "", err
	}
	_, err = file.Write(resp.Data)
	if err != nil {
		return "", err
	}
	file.Close()
	abs, err := filepath.Abs(file.Name())
	if err != nil {
		return "", err
	}
	return abs, nil
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
