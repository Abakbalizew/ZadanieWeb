package handlers

//
//
//

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"zadanieweb/databases"
	myerrors "zadanieweb/errors"
	"zadanieweb/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Главная страница
func Handle_api(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/api.html")
	if err != nil {
		panic(err)
	}

	userUUIDstring := models.AuthMap[models.AuthKey]
	if userUUIDstring != "" {
		cur_user := models.User{}

		userUUID, err := uuid.Parse(userUUIDstring)
		if err != nil {
			panic(err)
		}

		cur_user.FillUserWithUUID(userUUID)

		t.Execute(w, cur_user)
	} else {
		t.Execute(w, models.EmptyUser())
	}

}

// Обработчик страницы входа в аккаунт
func Handle_getLogin(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/getlogin.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, myerrors.Cur_error)
	myerrors.Cur_error.ErrMsg = ""
}

// Обработчик страницы входа в аккаунт
func Handle_getRegister(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/getregister.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, myerrors.Cur_error)
	myerrors.Cur_error.ErrMsg = ""
}

// Обработчик страницы с постами
func Handle_posts(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/posts.html")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", databases.ConnStr)
	if err != nil {
		fmt.Printf("Ошибка подключения бд :( ")
		panic(err)
	}
	defer db.Close()

	//Наш пользователь
	cur_user := models.User{}

	//UUID пользователя в виде строки
	userUUIDstring := models.AuthMap[models.AuthKey]
	if userUUIDstring != "" {
		cur_user = models.User{}

		//Переводим строку в тип uuid.UUID
		userUUID, err := uuid.Parse(userUUIDstring)
		if err != nil {
			panic(err)
		}
		//Заполняем нашего пользователя данными из бд
		cur_user.FillUserWithUUID(userUUID)

		//Если пользователь - читатель
		if cur_user.Role == "Reader" {
			//Выбираем из бд тольке те посты, которые опубликованы
			posts, err := db.Query("SELECT * FROM posts "+
				"WHERE status = $1 ORDER BY createdat DESC", "Published")
			if err != nil {
				fmt.Printf("Ошибка с поиском постов в бд :( ")
				panic(err)
			}
			defer posts.Close()

			for posts.Next() {
				post := models.Post{}
				err = posts.Scan(&post.PostUUID, &post.AuthorUUID,
					&post.IdempotencyKey, &post.Title, &post.Content,
					&post.CreatedAt, &post.UpdatedAt, &post.Status)
				if err != nil {
					panic(err)
				}
				if post.CreatedAt.After(post.UpdatedAt) {
					post.LastEditedAt = post.CreatedAt
				} else {
					post.LastEditedAt = post.UpdatedAt
				}
				//Посты уже отобраны, все со статусом "Published".
				cur_user.Posts = append(cur_user.Posts, post)
			}
		}
		//Если пользователь - автор
		if cur_user.Role == "Author" {
			//Выбираем из бд тольке посты пользователя
			posts, err := db.Query("SELECT * FROM posts "+
				"WHERE authorid = $1 ORDER BY createdat DESC", cur_user.UserUUID)
			if err != nil {
				fmt.Printf("Ошибка с поиском постов в бд :( ")
				panic(err)
			}
			defer posts.Close()

			for posts.Next() {
				post := models.Post{}
				err = posts.Scan(&post.PostUUID, &post.AuthorUUID,
					&post.IdempotencyKey, &post.Title, &post.Content,
					&post.CreatedAt, &post.UpdatedAt, &post.Status)
				if err != nil {
					panic(err)
				}
				if post.CreatedAt.After(post.UpdatedAt) {
					post.LastEditedAt = post.CreatedAt
				} else {
					post.LastEditedAt = post.UpdatedAt
				}

				//Мы уже отобрали только посты нашего пользователя.
				cur_user.Posts = append(cur_user.Posts, post)
			}
		}

	} else {
		cur_user = models.EmptyUser()
	}

	t.Execute(w, cur_user)

}

// Выход из аккаунта с главной страницы
func Handle_exit(w http.ResponseWriter, r *http.Request) {
	models.AuthMap[models.AuthKey] = ""

	http.Redirect(w, r, "/api", http.StatusSeeOther)
}

func Handle_exit_from_posts(w http.ResponseWriter, r *http.Request) {
	models.AuthMap[models.AuthKey] = ""

	http.Redirect(w, r, "/api/posts", http.StatusSeeOther)
}

// Обрабатывает страницу редактирования поста
func PostEditHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/postedit.html")
	if err != nil {
		panic(err)
	}

	vars := mux.Vars(r)
	postUUID := vars["uuid"]

	cur_post := models.Post{}

	db, err := sql.Open("postgres", databases.ConnStr)
	if err != nil {
		fmt.Printf("Ошибка подключения бд :( ")
		panic(err)
	}
	defer db.Close()

	err = db.QueryRow("SELECT * FROM posts WHERE postid = $1",
		postUUID).Scan(&cur_post.PostUUID, &cur_post.AuthorUUID, &cur_post.IdempotencyKey,
		&cur_post.Title, &cur_post.Content, &cur_post.CreatedAt, &cur_post.UpdatedAt, &cur_post.Status)
	if err != nil {
		fmt.Printf("Такого поста нет :( ")
		panic(err)
	}

	t.Execute(w, cur_post)

}

func PostDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postUUID := vars["uuid"]

	db, err := sql.Open("postgres", databases.ConnStr)
	if err != nil {
		fmt.Printf("Ошибка подключения бд :( ")
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM posts WHERE postid = $1", postUUID)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/api/posts", http.StatusSeeOther)
}
