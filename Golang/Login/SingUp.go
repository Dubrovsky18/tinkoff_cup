package Login

import (
	"database/sql"
	"fmt"
	"github.com/Dubrovsky18/tinkoff_cup/Login/Conn"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

// User структура, представляющая пользователя
type Company struct {
	Login    string
	Password string
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS company (
			login TEXT PRIMARY KEY,
			password TEXT NOT NULL
		);
	`)
	return err
}

// insertUser добавляет нового пользователя в таблицу users
func insertUser(db *sql.DB, company Company) error {
	_, err := db.Exec(`
		INSERT INTO company (login, password)
		VALUES ($1, $2);
	`, company.Login, company.Password)
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
	fmt.Println(err)
}
func RegistrationHandlerPost(w http.ResponseWriter, r *http.Request) {
	// Получаем данные пользователя из формы
	login := r.FormValue("login")
	password := r.FormValue("password")

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем пользователя
	company := Company{Login: login, Password: string(hashedPassword)}

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
	err = insertUser(db, company)
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
