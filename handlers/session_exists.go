package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"miauw.social/auth/database"
)

type ExistsUserSessionData struct {
	SID string
}

func ExistsUserSession(db *gorm.DB, rawData []byte) (Response, error) {
	var sessionData ExistsUserSessionData
	err := json.Unmarshal(rawData, &sessionData)
	if err != nil {
		return Response{
			Content: nil,
			Status: ResponseStatus{
				Code:   422,
				Title:  "Not processable!",
				Detail: "The data send to the worker was not processable.",
				Type:   "https://auth.miauw.social/session/exists/not-processable",
			},
		}, err
	}
	rdb := database.RedisConn()
	amount, err := rdb.Exists(context.Background(), sessionData.SID).Result()
	if err != nil {
		return Response{
			Content: nil,
			Status: ResponseStatus{
				Code:   400,
				Title:  "Error occured!",
				Detail: fmt.Sprintf("An unknown error occured: %v", err),
				Type:   "https://auth.miauw.social/session/exists/unknown-error",
			},
		}, err
	} else if err == redis.Nil {
		return Response{
			Content: nil,
			Status: ResponseStatus{
				Code:   404,
				Title:  "Session not found!",
				Detail: "This session for the user does not exist.",
				Type:   "https://auth.miauw.social/session/exists/not-found",
			},
		}, err
	}
	return Response{
		Content: true,
		Status: ResponseStatus{
			Code:   200,
			Detail: fmt.Sprintf("found: %v sessions", amount),
		},
	}, nil
}
