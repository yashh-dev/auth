package handlers

import (
	"encoding/json"

	"github.com/gofrs/uuid"
	"miauw.social/auth/database"
	"miauw.social/auth/database/models"
	"miauw.social/auth/security"
)

type UserCreateData struct {
	ID       string
	Password string
}

func UserCreate(rawData []byte) (Response, error) {
	var userCreateData UserCreateData
	err := json.Unmarshal(rawData, &userCreateData)
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
	passwordHash, err := security.EncryptPassword(userCreateData.Password)
	db := database.Conn()
	account := models.Account{PasswordHash: passwordHash, Base: models.Base{ID: uuid.FromStringOrNil(userCreateData.ID)}}
	db.Create(&account)
	vid, _ := security.GenerateJWT(account.ID.String())
	return Response{
		Content: vid,
		Status: ResponseStatus{
			Code: 201,
		},
	}, err
}
