package handlers

import (
	"encoding/json"
	"fmt"

	"miauw.social/auth/database"
	"miauw.social/auth/database/models"
	"miauw.social/auth/security"
)

type UserVerifyData struct {
	Token string
}

func UserVerify(rawData []byte) (Response, error) {
	var userVerifyData UserVerifyData
	err := json.Unmarshal(rawData, &userVerifyData)
	if err != nil {
		return Response{
			Content: nil,
			Status: ResponseStatus{
				Code:   422,
				Title:  "Not processable!",
				Detail: "The data send to the worker was not processable.",
				Type:   "https://auth.miauw.social/verify/not-processable",
			},
		}, err
	}
	claims, err := security.VerifyJWT(userVerifyData.Token)
	if err != nil {
		return Response{
			Content: nil,
			Status: ResponseStatus{
				Code:   403,
				Title:  "Token not validated!",
				Detail: fmt.Sprintf("Your token could not be validated. %v", err),
				Type:   "https://auth.miauw.social/verify/token-invalid",
			},
		}, err
	}
	db := database.Conn()
	var account models.Account
	db.Where("id::text = ?", claims["sub"]).Find(&account)
	return Response{
		Content: true,
		Status: ResponseStatus{
			Code: 202,
		},
	}, err
}
