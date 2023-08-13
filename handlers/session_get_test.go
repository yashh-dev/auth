package handlers

import (
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestSessionGet(t *testing.T){

	validaSessionData := GetUserSessionData{
		ID: "894cf31c-cb98-45f0-bce8-63848e09689e",
	}

	validaSessionDataRaw ,err:=json.Marshal(validaSessionData)
	if err != nil{
		t.Error(err)
	}
	
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
	
	t.Run("get user session",func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM `sessions` WHERE \"user\"(.+)text  = (.+) AND created_at >= (.+) LIMIT 1").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("894cf31c-cb98-45f0-bce8-63848e09689e"))
		
		resoponse,err:=GetUserSession(gormDb,validaSessionDataRaw)

		assert.Empty(t,err)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		assert.Equal(t,200,resoponse.Status.Code)
	})
}