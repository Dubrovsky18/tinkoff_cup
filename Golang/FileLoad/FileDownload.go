ackage FileLoad

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	// Проверяем, залогинен ли пользователь
	session, err := r.Cookie("session")
	if err != nil || session.Value == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("templates/download.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "download", nil)

	done := make(chan bool)
	go FileDownload(w, r, done)
	select {
	case <-done:
		// Файл успешно отправлен
	case <-time.After(time.Second * 50):
		// Прошло 10 секунд, но файл так и не пришел
		http.Error(w, "Время ожидания истекло", http.StatusInternalServerError)
	}
}

func FileDownload(w http.ResponseWriter, r *http.Request, done chan<- bool) {
	conn, err := net.Dial("tcp", "localhost:4999")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		done <- true
		return
	}
	defer conn.Close()

	// Получаем имя файла от сервера
	filename := make([]byte, 1024)
	n, err := conn.Read(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		done <- true
		return
	}

	// Удаляем лишние символы из имени файла
	filename = filename[:n]
	filename = bytes.Trim(filename, "\x00")

	// Открываем файл и проверяем на ошибки
	file, err := os.Open(fmt.Sprintf("FileLoad/FilesWebSiteOut/%s", filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		done <- true
		return
	}
	defer file.Close()

	// Получаем информацию о файле
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		done <- true
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
		done <- true
		return
	}

	done <- true
}
