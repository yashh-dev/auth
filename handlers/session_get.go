package handlers

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type GetUserSessionData struct {
	UserID string
}

func GetUserSession(db *gorm.DB, rawData []byte) (Response, error) {
	var sessionData GetUserSessionData
	err := json.Unmarshal(rawData, &sessionData)
	if err != nil {
		return Response{
			Content: nil,
			Status: ResponseStatus{
				Code:   422,
				Title:  "Not processable!",
				Detail: "The data send to the worker was not processable.",
				Type:   "https://auth.miauw.social/login/not-processable",
			},
		}, err
	}
	// var sessions []models.Session
	results := make(map[string]interface{})
	db.Where("user  = ?", sessionData.UserID).Take(&results)
	fmt.Println(results)
	return Response{
		Content: &results,
		Status: ResponseStatus{
			Code: 200,
		},
	}, nil
}
