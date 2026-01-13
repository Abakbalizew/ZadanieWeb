package main

import (
	"net/http"
	"zadanieweb/handlers"
	"zadanieweb/jwttokens"
	"zadanieweb/models"

	"github.com/gorilla/mux"

	_ "zadanieweb/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Zadan1e
// @description Здесь можно делиться своими мыслями и искать себе кумиров.
// @version 1.0

// @contact.name   S4pan
// @contact.email  z4har.sapan@gmail.com

// @host localhost:8080
// @BasePath /api

func main() {
	//Обнуляем ошибки, чтобы они не выводились на экранах регистрации и входа в аккаунт просто так
	models.AuthMap[models.AuthKey] = ""

	//Маршрутизатор
	rtr := mux.NewRouter()

	//API документация
	rtr.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	//Главная страница
	rtr.HandleFunc("/api", handlers.MainPageHandler).Methods("GET")
	//Страница входа в аккаунт
	rtr.HandleFunc("/api/auth/loginPage", handlers.LoginPageHandler).Methods("GET")
	//Страница регистрации
	rtr.HandleFunc("/api/auth/registerPage", handlers.RegisterPageHandler).Methods("GET")
	//Выход из аккаунта
	rtr.HandleFunc("/api/posts/exitFromMainPage", handlers.ExitFromMainPageHandler).Methods("GET")
	//Выход из аккаунта со страницы постов
	rtr.HandleFunc("/api/posts/exitFromPosts", handlers.ExitFromPostsHandler).Methods("GET")

	//Посты
	rtr.HandleFunc("/api/posts", handlers.PostsHandler).Methods("GET")
	//Создание (с проверкой jwt токена)
	rtr.HandleFunc("/api/posts/create{uuid:(?:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})?}",
		jwttokens.CheckTokenMiddleWare(handlers.HandlePostCreationORSavingToDraft("Published"))).Methods("POST")
	//Сохранение в черновик (с проверкой jwt токена)
	rtr.HandleFunc("/api/posts/save{uuid:(?:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})?}",
		jwttokens.CheckTokenMiddleWare(handlers.HandlePostCreationORSavingToDraft("Draft"))).Methods("POST")
	//Редактирование отдельного поста (с проверкой jwt токена)
	rtr.HandleFunc("/api/posts/{uuid:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}}",
		jwttokens.CheckTokenMiddleWare(handlers.PostEditHandler)).Methods("GET")
	//Удаление поста (с проверкой jwt токена)
	rtr.HandleFunc("/api/posts/delete{uuid:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}}",
		jwttokens.CheckTokenMiddleWare(handlers.PostDeleteHandler)).Methods("POST")

	//Пост-запрос на вход в учётную запись
	rtr.HandleFunc("/api/auth/login", handlers.LoginRequestHandler).Methods("POST")
	//Пост-запрос на регистрацию
	rtr.HandleFunc("/api/auth/register", handlers.RegisterRequestHandler).Methods("POST")

	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)
}
