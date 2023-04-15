package FileUpload

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

func HandleMain(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/main.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "main", nil)

}

func FileUpload(w http.ResponseWriter, r *http.Request) {

	link := r.FormValue("link")
	fmt.Printf("Link:", link)

	// Сохраняем загруженный файл на сервере
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Создаем новый файл на сервере и копируем содержимое загруженного файла в него
	out, err := os.Create(header.Filename)
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

	// Выводим сообщение об успешной загрузке файла на странице
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("File %s uploaded successfully", header.Filename)))

	http.Redirect(w, r, "/", http.StatusSeeOther)

}
