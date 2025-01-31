package data

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestUserOperations(t *testing.T) {
	// Initialize test database connection
	dbURI := os.Getenv("DATABASE_URL")
	if dbURI == "" {
		t.Skip("DATABASE_URL not set")
	}
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		t.Error("fail to conect database")
	}
	data := &Data{
		Db: db,
	}
	fakeUser := &User{
		Email: "jdr666@foxmail.com",
	}
	err = data.CreateUser(fakeUser)
	if err != nil {
		t.Errorf("fail to create user: %s", err.Error())
	}
	if isUser, _ := data.CheckUser(fakeUser.Email); isUser != true {
		t.Error("fail to check user")
	}

	notifications, err := data.GetUserAllNotifications(fakeUser.UserID)
	if err != nil {
		t.Error("fail to get user notifications")
	}
	if len(notifications) != 0 {
		t.Error("notifications should be empty")
	}

	email, err := data.GetUserEmail(fakeUser.UserID)
	if err != nil {
		t.Error("fail to get user email")
	}
	if email != fakeUser.Email {
		t.Error("email should be equal")
	}
	err = data.DeleteUser(fakeUser.UserID)
	if err != nil {
		t.Error("fail to delete user")
	}

	noUser, err := data.CheckUser(fakeUser.Email)
	if err != nil {
		t.Error("fail to check user")
	}
	if noUser != false {
		t.Error("user should not exist")
	}

}
