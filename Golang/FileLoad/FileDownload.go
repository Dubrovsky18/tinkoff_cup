package FileLoad

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	// Проверяем, залогинен ли пользователь
	session, err := r.Cookie("session")
	if err != nil || session.Value == "" || session.Value != user.Name {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("templates/download.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "download", nil)

}

func HandleFileDownload(w http.ResponseWriter, r *http.Request) {
	// Проверяем существование файла

	filePathOut = fmt.Sprintf("Test/%s/chrome-%s-video.mp4", user.Name, user.Name)
	_, err := os.Stat(filePathOut)
	if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Открываем файл и проверяем на ошибки
	file, err := os.Open(filePathOut)
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
	w.Header().Set("Content-Length", string(fileInfo.Size()))

	// Копируем содержимое файла в ResponseWriter
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
