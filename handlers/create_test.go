package handlers

import (
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)




 func TestCreate(t *testing.T){

	validaUserData := UserCreateData{
		ID: "894cf31c-cb98-45f0-bce8-63848e09689e",
		Password: "test@123",
	}

	validaUserDataRaw ,err:=json.Marshal(validaUserData)
	if err != nil{
		t.Error(err)
	}
	invalidUserDataRaw := []byte(`id:894cf31c-cb98-45f0-bce8-63848e09689e,password:test@123`)
	

	db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error on '%s' stub database connection", err)
		}

		gormDb,err:=gorm.Open(mysql.New(mysql.Config{
			Conn: db,
			SkipInitializeWithVersion: true,
		}))

		if err != nil{
			t.Error("gorm.Err",err)
		}
	
	t.Run("user data inserts into the db and token recieved",func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `accounts` .`id`,`deleted_at`,`password_hash`,`verified`. VALUES (.+)").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		response,err:=UserCreate(gormDb,validaUserDataRaw)

		assert.Empty(t,err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.Equal(t,201,response.Status.Code)
		assert.NotEmpty(t,response.Content)
	})
	
	t.Run("invalid user data fails with code 422",func(t *testing.T) {
		response,err:=UserCreate(gormDb,invalidUserDataRaw)

		if err != nil{
			t.Log("handler.create",err)
			assert.Equal(t,422,response.Status.Code)
		}
	})
}