package Login

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Dubrovsky18/tinkoff_cup/Login/Conn"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// User структура, представляющая пользователя
type User struct {
	Login    string
	Password string
	Team     string
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user_work (
			login TEXT PRIMARY KEY,
			password TEXT NOT NULL,
			team TEXT 
		);
	`)
	return err
}

// insertUser добавляет нового пользователя в таблицу users
func insertUser(db *sql.DB, company User) error {
	_, err := db.Exec(`
		INSERT INTO user_work (login, password,team)
		VALUES ($1, $2, $3);
	`, company.Login, company.Password, company.Team)
	return err
}

// handleRegister обработчик для страницы регистрации
func RegistrationHandlerGet(w http.ResponseWriter, r *http.Request) {
	// Отображаем форму регистрации
	t, err := template.ParseFiles("templates/registration.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	err = t.ExecuteTemplate(w, "registration", nil)

}
func RegistrationHandlerPost(w http.ResponseWriter, r *http.Request) {
	// Получаем данные пользователя из формы
	login := r.FormValue("login")
	password := r.FormValue("password")
	team := r.FormValue("team")

	if team == "" {
		team = "default"
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем пользователя
	user := User{Login: login, Password: string(hashedPassword), Team: team}

	db, err := Conn.OpenDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Создаем таблицу users, если ее нет
	err = createTable(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Добавляем пользователя в базу данных PostgreSQL
	err = insertUser(db, user)
	if err != nil {
		// Если пользователь уже существует, то выведем ошибку
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Перенаправляем пользователя на страницу после успешной регистрации
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
