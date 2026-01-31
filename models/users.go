package models

///Модель пользователя, а также методы и функции для работы с этой моделью

import (
	"database/sql"
	"fmt"
	"zadanieweb/databases"

	"github.com/google/uuid"
)

//AuthMap - состояние, хранящее в себе uuid аутентифицированного пользователя
//Если по ключу AuthKey ничего не лежит, то пользователь не аутентифицирован

// Состояние хранит в себе uuid строку аутентифицированного пользователя
var AuthMap = make(map[string]string)

// Ключ, по которому uuid будет храниться в состоянии
var AuthKey = "UserAuthorizationUUID"

// Модель пользователя
type User struct {
	UserUUID       uuid.UUID
	Email          string
	Role           string
	AccessToken    string
	RefreshToken   string
	HashedPassword []byte
	//Слайс постов, которые должны отображаться данному пользователю
	Posts []Post
	//Сообщение об ошибке, если такое нужно вывести
	ErrMsg string
}

// Метод, возвращающий пустого пользователя
func EmptyUser() User {
	return User{
		Email: "",
		Role:  "",
	}
}

// Заполняет переменную типа User
// соответствующими коду uuid данными из базы данных.
func (cur_user *User) FillUserWithUUID(userUUID uuid.UUID) {
	db, err := sql.Open("postgres", databases.ConnStr)
	if err != nil {
		fmt.Printf("Ошибка: models/users.go/FillUserWithUUID - 1")
		panic(err)
	}
	defer db.Close()

	err = db.QueryRow("SELECT userid, email, passwordhash, role_, accesstoken, refreshtoken FROM users WHERE userid = $1",
		userUUID).Scan(&cur_user.UserUUID, &cur_user.Email, &cur_user.HashedPassword,
		&cur_user.Role, &cur_user.AccessToken, &cur_user.RefreshToken)
	if err != nil {
		fmt.Printf("Ошибка: models/users.go/FillUserWithUUID - 2")
		panic(err)
	}
}
