package handlers

///Обработчики get-запросов

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

// MainPageHandler Главная страница
// @Summary Главная страница
// @Router /api [get]
// @Tags Main
func MainPageHandler(w http.ResponseWriter, r *http.Request) {
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

// LoginPageHandler Страница входа в аккаунт
// @Summary Страница входа в аккаунт
// @Description Имеются 2 формы для почты и пароля.
// @Router /api/auth/loginPage [get]
// @Tags Authoriztion
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/getlogin.html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, myerrors.Cur_error)

	myerrors.Cur_error.ErrMsg = ""
}

// RegisterPageHandler Страница регистрации
// @Summary Страница регистрации
// @Description Имеются 3 формы для почты, пароля и роли.
// @Router /api/auth/registerPage [get]
// @Tags Authoriztion
func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/getregister.html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, myerrors.Cur_error)

	myerrors.Cur_error.ErrMsg = ""
}

// PostsHandler Страница постов
// @Summary Страница постов
// @Description Страница постов отображается по разному в зависимости от роли и наличия авторизации.
// @Router /api/posts [get]
// @Tags Posts
func PostsHandler(w http.ResponseWriter, r *http.Request) {
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

				//Проверяем, есть ли у поста изображение:
				imageRows, err := db.Query("SELECT imageurl FROM images WHERE postid = $1", post.PostUUID)
				if err != nil {
					fmt.Printf("Ошибка поиска значений из бд! :( ")
					panic(err)
				} //Если нет:
				if !imageRows.Next() {
					post.ImageUrl = ""
				} else { //Если да:
					imageRows.Scan(&post.ImageUrl)
				}

				//Получаем email автора
				err = db.QueryRow("SELECT email FROM users WHERE userid = $1", post.AuthorUUID).Scan(&post.AuthorEmail)
				if err != nil {
					panic(err)
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

				//Проверяем, есть ли у поста изображение:
				imageRows, err := db.Query("SELECT imageurl FROM images WHERE postid = $1", post.PostUUID)
				if err != nil {
					fmt.Printf("Ошибка поиска значений из бд! :( ")
					panic(err)
				} //Если нет:
				if !imageRows.Next() {
					post.ImageUrl = ""
				} else { //Если да:
					imageRows.Scan(&post.ImageUrl)
				}

				//Мы уже отобрали только посты нашего пользователя.
				cur_user.Posts = append(cur_user.Posts, post)
			}
		}
		//Если пользователь не авторизован
	} else {
		cur_user = models.EmptyUser()
	}

	t.Execute(w, cur_user)

}

// ExitFromMainPageHandler Выход из аккаунта с главной страницы
// @Summary Выход из аккаунта с главной страницы
// @Description Выход из аккаунта, мы остаёмся на той же странице
// @Router /api/posts/exitFromMainPage [get]
// @Tags Authoriztion
func ExitFromMainPageHandler(w http.ResponseWriter, r *http.Request) {
	models.AuthMap[models.AuthKey] = ""

	http.Redirect(w, r, "/api", http.StatusSeeOther)
}

// ExitFromMainPageHandler Выход из аккаунта со страницы постов
// @Summary Выход из аккаунта со страницы постов
// @Description Выход из аккаунта, мы остаёмся на той же странице.
// @Router /api/posts/exitFromPosts [get]
// @Tags Authoriztion
func ExitFromPostsHandler(w http.ResponseWriter, r *http.Request) {
	models.AuthMap[models.AuthKey] = ""

	http.Redirect(w, r, "/api/posts", http.StatusSeeOther)
}

// PostEditHandler Редактирование поста
// @Summary Редактирование поста
// @Description Страница редактирования конкретного поста.
// @param uuid path string false "Post uuid"
// @Router /api/posts/{uuid} [get]
// @Tags Posts
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

	//Проверим, есть ли изображение:
	imageRows, err := db.Query("SELECT imageurl FROM images WHERE postid = $1", postUUID)
	if err != nil {
		fmt.Printf("Ошибка поиска значений из бд! :( ")
		panic(err)
	}
	if imageRows.Next() {
		imageRows.Scan(&cur_post.ImageUrl)
	}

	t.Execute(w, cur_post)
}
