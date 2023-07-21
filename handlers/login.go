package handlers

import (
	"encoding/json"

	"gorm.io/gorm"
	"miauw.social/auth/database/models"
	"miauw.social/auth/security"
)

type UserSessionResponse struct {
	SID string
}

type UserLoginData struct {
	ID       string
	Password string
}

func UserLogin(db *gorm.DB, rawData []byte) (Response, error) {
	var userLoginData UserLoginData
	err := json.Unmarshal(rawData, &userLoginData)
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
	var account models.Account
	db.Where("id = ?", userLoginData.ID).First(&account)
	if !account.Verified {
		return Response{
			Content: nil,
			Status: ResponseStatus{
				Code:   403,
				Title:  "Account not verified",
				Detail: "Your account is not verified. Please verify it first.",
				Type:   "https://auth.miauw.social/login/not-verified",
			},
		}, err
	}
	ok, err := security.VerifyPassword(account.PasswordHash, userLoginData.Password)
	if !ok {
		return Response{
			Content: nil,
			Status: ResponseStatus{
				Code:   401,
				Title:  "Wrong password",
				Detail: "The submitted password does not match the password in the database.",
				Type:   "https://auth.miauw.social/login/wrong-password",
			},
		}, err
	}
	session := models.Session{UserID: account.ID}
	db.Create(&session)
	return Response{
		Content: UserSessionResponse{SID: session.ID.String()},
		Status: ResponseStatus{
			Code: 200,
		},
	}, err
}
