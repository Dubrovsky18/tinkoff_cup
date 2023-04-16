package Login

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Dubrovsky18/tinkoff_cup/Login/Conn"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
)

// User структура, представляющая пользователя
type Team struct {
	name string
}

func createTeamTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS teams (
			name TEXT PRIMARY KEY
		);
	`)
	return err
}

// insertUser добавляет новую группу в базу
func insertGroup(db *sql.DB, team Team) error {
	_, err := db.Exec(`
		INSERT INTO teams (name)
		VALUES ($1);
	`, team.name)
	return err
}

// handleRegister обработчик для страницы регистрации
func RegistrationTeamHandlerGet(w http.ResponseWriter, r *http.Request) {
	// Отображаем форму регистрации
	t, err := template.ParseFiles("templates/registration_team.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	err = t.ExecuteTemplate(w, "registration", nil)
	fmt.Println(err)
}
func RegistrationTeamHandlerPost(w http.ResponseWriter, r *http.Request) {
	// Получаем данные группы из формы
	name := r.FormValue("team")

	// Создаем пользователя
	team := Team{name: name}

	db, err := Conn.OpenDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Создаем таблицу users, если ее нет
	err = createTeamTable(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Добавляем группу в базу данных PostgreSQL
	err = insertGroup(db, team)
	if err != nil {
		// Если группа уже существует, то выведем ошибку
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			http.Error(w, "Team already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Перенаправляем пользователя на страницу после успешной регистрации
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
