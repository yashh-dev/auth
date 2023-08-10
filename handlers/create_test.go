package handlers

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCreate(t *testing.T){
		rawData := []byte(`{"id":"945123486","password":"test@123"}`)
	
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
	
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `accounts` .`id`,`deleted_at`,`password_hash`,`verified`. VALUES (.+)").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		response,err:=UserCreate(gormDb,rawData)

		if err != nil{
			t.Error("handler.create",err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.Equal(t,201,response.Status.Code)
}