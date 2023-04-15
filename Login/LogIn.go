package Login

import (
	"database/sql"
	"fmt"
	"github.com/Dubrovsky18/tinkoff_cup/Login/Conn"
	_ "github.com/Dubrovsky18/tinkoff_cup/Login/Conn"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

func LoginHandleGet(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "login", nil)

}

func LoginHandlePost(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	db, err := Conn.OpenDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Ищем пользователя в базе данных PostgreSQL
	var company Company
	err = db.QueryRow("SELECT * FROM company WHERE login=$1", login).Scan(&company.Login, &company.Password)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid login or password", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(company.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid login or password", http.StatusBadRequest)
		return
	}

	// Устанавливаем cookie с информацией о пользователе
	cookie := http.Cookie{
		Name:     "user",
		Value:    company.Login,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	// Перенаправляем пользователя на главную страницу
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}
