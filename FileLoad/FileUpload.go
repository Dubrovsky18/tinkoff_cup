package FileLoad

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/upload.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "upload", nil)
}
func FileUpload(w http.ResponseWriter, r *http.Request) {
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
	out, err := os.Create(filepath.Join("FileLoad/FilesWebSiteIn", header.Filename))
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
