package FileLoad

import (
	"fmt"
	_ "github.com/gorilla/sessions"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
)

type User struct {
	session  string
	nickname string
}

var userName string

func HandleDownload(w http.ResponseWriter, r *http.Request) {

	// Проверяем, залогинен ли пользователь
	session, err := r.Cookie("session")
	if err != nil || session.Value == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userName = session.Value

	t, err := template.ParseFiles("templates/download.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "download", nil)

	http.HandleFunc("/downloadFile", FileDownload)
	http.ListenAndServe(":4999", nil)
	fmt.Fprintf(w, session.Value, "Download page")

}

func FileDownload(w http.ResponseWriter, r *http.Request) {

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Не указано имя файла", http.StatusBadRequest)
		return
	}

	// Получаем путь к файлу из параметра запроса
	filePath := fmt.Sprintf("FileLoad/FilesWebSiteOut/tests/%s/%s", userName, "OutFile")

	// Открываем файл и проверяем на ошибки
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Получаем информацию о файле
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки
	w.Header().Set("Content-Disposition", "attachment; filename="+fileInfo.Name())
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// Копируем содержимое файла в ResponseWriter
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
