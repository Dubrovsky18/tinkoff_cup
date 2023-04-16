package FileLoad

import (
	"fmt"
	_ "github.com/gorilla/sessions"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что пользователь аутентифицирован
	session, err := r.Cookie("session")
	if err != nil || session.Value == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("templates/upload.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.ExecuteTemplate(w, "upload", nil)

}

func FileUpload(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что пользователь аутентифицирован
	session, err := r.Cookie("session")
	if err != nil || session.Value == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Получаем информацию о загружаемом файле
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	serverPath := strings.ReplaceAll((session.Value + "_" + fileHeader.Filename), " ", "")

	// Сохраняем загруженный файл на сервере
	filePath := fmt.Sprintf("FileLoad/FilesWebSiteIn/%s", serverPath)
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем файл на другой сервер для обработки
	f, err := os.Open(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	resp, err := http.Post("http://localhost:5000/run-tests?filename="+serverPath, "application/octet-stream", f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Перенаправляем пользователя на страницу с загрузкой
	http.Redirect(w, r, "/download", http.StatusSeeOther)
}
