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

	go HandleFileDownload(w, r)
	t.ExecuteTemplate(w, "download", nil)

}

func HandleFileDownload(w http.ResponseWriter, r *http.Request) {
	// Проверяем существование файла

	filePathOut = "test.py"
	filePathOut = fmt.Sprintf("%s/%s", user.Name, "test.py")
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
	w.Header().Set("Content-Length", string(99))

	// Копируем содержимое файла в ResponseWriter
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
