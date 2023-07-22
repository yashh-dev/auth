package handlers

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"miauw.social/auth/database/models"
)

type GetUserSessionData struct {
	ID string
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
	deltaTtl := time.Now().Add(-12 * time.Hour)
	var sessions []models.Session
	db.Where(`"user"::text  = ?`, sessionData.ID).Where("created_at >= ?", deltaTtl).Take(&sessions)
	return Response{
		Content: &sessions,
		Status: ResponseStatus{
			Code: 200,
		},
	}, nil
}
