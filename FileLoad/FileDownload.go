package FileLoad

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
)

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/download.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "download", nil)

}

func FileDownload(w http.ResponseWriter, r *http.Request) {

	// Получаем путь к файлу из параметра запроса
	filePath := "FileLoad/FilesWebSiteOut/Whale-52.txt"

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

//func FileDownload(w http.ResponseWriter, r *http.Request) {
//	// Путь к файлу, который нужно вернуть
//	filePath := fmt.Sprintf("FileWebSiteOut/%s", "Whale-52.txt")
//
//	// Открываем файл и проверяем на ошибки
//	file, err := os.Open(filePath)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer file.Close()
//
//	// Получаем информацию о файле
//	fileInfo, err := file.Stat()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	// Устанавливаем заголовки
//	w.Header().Set("Content-Disposition", "attachment; filename="+fileInfo.Name())
//	w.Header().Set("Content-Type", "application/octet-stream")
//	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
//
//	// Копируем содержимое файла в ResponseWriter
//	_, err = io.Copy(w, file)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}
