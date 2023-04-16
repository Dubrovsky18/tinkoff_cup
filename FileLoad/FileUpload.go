package FileLoad

import (
	"fmt"
	_ "github.com/gorilla/sessions"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "upload", nil)

	fmt.Fprintf(w, session.Value, "Upload page")

}
func FileUpload(w http.ResponseWriter, r *http.Request) {

	session, err := r.Cookie("session")

	if err != nil || session.Value == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Продолжаем загрузку файла
	link := r.FormValue("link")
	fmt.Printf("Link: %s", link)
	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := filepath.Base(fileHeader.Filename)
	fmt.Printf("\n File Name: %s", filename)

	// Сохраняем загруженный файл на сервере
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Создаем новый файл на сервере и копируем содержимое загруженного файла в него
	out, err := os.Create(filepath.Join(fmt.Sprintf("FileLoad/FilesWebSiteIn/tests/%s/%s", session.Value, header.Filename)))
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

	http.Redirect(w, r, "/download", http.StatusSeeOther)

}
