package models

import (
	"database/sql"
	"fmt"
	"zadanieweb/databases"

	"github.com/google/uuid"
)

// Состояние хранит в себе uuid строку аутентифицированного пользователя
var AuthMap = make(map[string]string)

// Ключ, по которому uuid будет храниться в состоянии
var AuthKey = "UserAuthorizationUUID"

type User struct {
	UserUUID       uuid.UUID
	Email          string
	Role           string
	AccessToken    string
	RefreshToken   string
	HashedPassword []byte
}

func EmptyUser() User {
	return User{
		Email: "",
		Role:  "",
	}
}

func (cur_user *User) FillUserWithUUID(userUUID uuid.UUID) {
	db, err := sql.Open("postgres", databases.ConnStr)
	if err != nil {
		fmt.Printf("Ошибка подключения бд :( ")
		panic(err)
	}
	defer db.Close()

	err = db.QueryRow("SELECT * FROM users WHERE userid = $1",
		userUUID).Scan(&cur_user.UserUUID, &cur_user.Email, &cur_user.Role,
		&cur_user.RefreshToken, &cur_user.AccessToken, &cur_user.HashedPassword)
	if err != nil {
		fmt.Printf("Ошибка сохранения значений в бд")
		panic(err)
	}
}
