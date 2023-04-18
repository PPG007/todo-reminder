package gocq

import (
	"context"
	"errors"
	"fmt"
	"github.com/qiniu/qmgo"
	"strings"
	"todo-reminder/log"
	"todo-reminder/model"
)

const (
	CMD_APPROVE = "/approve"
	CMD_REJECT  = "/reject"
	CMD_BLOCK   = "/block"

	SUB_CMD_FRIEND = "friend"
	SUB_CMD_OPENAI = "openAI"
)

type Command struct {
	Name string
	Args []string
}

func ParseCommand(content string) (Command, error) {
	if !strings.HasPrefix(content, "/") {
		return Command{}, errors.New("invalid command")
	}
	args := strings.Split(content, " ")
	cmd := Command{
		Name: strings.TrimSpace(args[0]),
	}
	for i, str := range args {
		if i == 0 {
			continue
		}
		cmd.Args = append(cmd.Args, strings.TrimSpace(str))
	}
	return cmd, nil
}

func (c *Command) Run(ctx context.Context) (string, error) {
	switch c.Name {
	case CMD_APPROVE:
		return c.runAsApprove(ctx)
	case CMD_REJECT:
		return c.runAsReject(ctx)
	case CMD_BLOCK:
		return c.runAsBlock(ctx)
	default:
		return "", errors.New("invalid command")
	}
}

func (c *Command) runAsApprove(ctx context.Context) (string, error) {
	if len(c.Args) != 2 {
		return "", errors.New("invalid arguments")
	}
	switch c.Args[0] {
	case SUB_CMD_FRIEND:
		_, err := model.CWaitingUser.GetWaitingOne(ctx, c.Args[1])
		if err != nil {
			if err == qmgo.ErrNoSuchDocuments {
				return "", errors.New("friend request from this user not found")
			}
			return "", err
		}
		err = goCqWs.handleFriendRequest(ctx, c.Args[1], true)
		if err != nil {
			return "", err
		}
	case SUB_CMD_OPENAI:
		user, err := model.CUser.ApproveOpenAI(ctx, c.Args[1])
		if err != nil {
			return "", err
		}
		sendNotifyToTargetUser(ctx, "approved", SUB_CMD_OPENAI, user.UserId)
	default:
		return "", errors.New("invalid arguments")
	}
	return "OK", nil
}

func (c *Command) runAsReject(ctx context.Context) (string, error) {
	if len(c.Args) != 2 {
		return "", errors.New("invalid arguments")
	}
	switch c.Args[0] {
	case SUB_CMD_FRIEND:
		_, err := model.CWaitingUser.GetWaitingOne(ctx, c.Args[1])
		if err != nil {
			if err == qmgo.ErrNoSuchDocuments {
				return "", errors.New("friend request from this user not found")
			}
			return "", err
		}
		err = goCqWs.handleFriendRequest(ctx, c.Args[1], false)
		if err != nil {
			return "", err
		}
	case SUB_CMD_OPENAI:
		user, err := model.CUser.GetByUserId(ctx, c.Args[1])
		if err != nil {
			return "", err
		}
		sendNotifyToTargetUser(ctx, "rejected", SUB_CMD_OPENAI, user.UserId)
	default:
		return "", errors.New("invalid arguments")
	}
	return "OK", nil
}

func (c *Command) runAsBlock(ctx context.Context) (string, error) {
	log.Warn("runAsBlock", map[string]interface{}{
		"cmd": c,
	})
	if len(c.Args) != 2 {
		return "", errors.New("invalid arguments")
	}
	switch c.Args[0] {
	case SUB_CMD_FRIEND:
		err := goCqWs.deleteFriend(ctx, c.Args[1])
		if err != nil {
			return "", err
		}
	case SUB_CMD_OPENAI:
		user, err := model.CUser.BlockOpenAI(ctx, c.Args[1])
		if err != nil {
			return "", err
		}
		sendNotifyToTargetUser(ctx, "blocked", SUB_CMD_OPENAI, user.UserId)
	default:
		return "", errors.New("invalid arguments")
	}
	return "OK", nil
}

func IsCommand(content string) bool {
	return strings.Contains(content, "/")
}

func sendNotifyToTargetUser(ctx context.Context, cmd, source, userId string) error {
	return goCqWs.SendPrivateStringMessage(ctx, fmt.Sprintf("Your are %s to use %s.", cmd, source), userId)
}
