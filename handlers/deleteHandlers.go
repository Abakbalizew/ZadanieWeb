package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"zadanieweb/databases"

	"github.com/gorilla/mux"
)

// PostDeleteHandler Удаление поста
// @Summary Удаление поста
// @Description Удаление конкретного поста.
// @param uuid path string false "Post uuid"
// @Router /api/posts/delete{uuid} [delete]
// @Tags Posts
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
