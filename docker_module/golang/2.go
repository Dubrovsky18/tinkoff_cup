package main

import (
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path/filepath"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/client"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.POST("/run-tests", runContainer)
    r.Run(":5000")
}

func runContainer(c *gin.Context) {
    // Аргументы
    filename := c.Query("filename")
    tempDir := c.Query("tempdir")

    // Конфигурация Docker
    dockerClient, err := client.NewEnvClient()
    if err != nil {
        log.Fatalf("Failed to create Docker client: %s", err.Error())
    }

    // Сохраняем в переменной путь до файла, tempDir - папка в которую бэк загрузил файл
    testFilePath := filepath.Join(tempDir, filename)

    // Создание Dockerfile из шаблона
    dockerfileTemplatePath := "Dockerfile-template"
    pathToTestInContainer := "/app/" + filename

    // Добавление в шаблон кастомных строк
    dockerfileCopyPath := filepath.Join(tempDir, "Dockerfile")
    dockerfileContent, err := ioutil.ReadFile(dockerfileTemplatePath)
    if err != nil {
        log.Fatalf("Failed to read Dockerfile template: %s", err.Error())
    }
    dockerfileContent += []byte("\nCOPY " + filename + " /app/" + filename + "\n")
    dockerfileContent += []byte("\nCMD [\"python\", \"" + pathToTestInContainer + "\", \" >> result\"]\n")

    err = ioutil.WriteFile(dockerfileCopyPath, dockerfileContent, 0644)
    if err != nil {
        log.Fatalf("Failed to create Dockerfile: %s", err.Error())
    }

    // Создаем контейнер и запускаем тест
    ctx := c.Request.Context()
    buildCtx, err := os.Open(tempDir)
    if err != nil {
        log.Fatalf("Failed to open build context: %s", err.Error())
    }
    defer buildCtx.Close()

    buildOptions := types.ImageBuildOptions{
        Dockerfile: "Dockerfile",
        Context:    buildCtx,
        Remove:     true,
    }

    buildResponse, err := dockerClient.ImageBuild(ctx, buildCtx, buildOptions)
    if err != nil {
        log.Fatalf("Failed to build Docker image: %s", err.Error())
    }
    defer buildResponse.Body.Close()

    scanner := bufio.NewScanner(buildResponse.Body)
    for scanner.Scan() {
        log.Println(scanner.Text())
    }

    containerConfig := &container.Config{
        Image:        buildResponse.ID,
        AttachStdout: true,
        AttachStderr: true,
        Cmd:          []string{"python", pathToTestInContainer},
    }
    hostConfig := &container.HostConfig{}

    containerResponse, err := dockerClient.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
    if err != nil {
        log.Fatalf("Failed to create Docker container: %s", err.Error())
    }

    err = dockerClient.ContainerStart(ctx, containerResponse.ID, types.ContainerStartOptions{})
    if err != nil {
        log.Fatalf("Failed to start Docker container: %s", err.Error())
    }

    containerLogsOptions := types.ContainerLogsOptions{
        ShowStdout: true,
        ShowStderr: true,
        Follow: true,
    }

    // Получение логов контейнера
    containerLogsResponse, err := dockerClient.ContainerLogs(ctx, containerResponse.ID, containerLogsOptions)
    if err != nil {
        log.Fatalf("Failed to retrieve container logs: %s", err.Error())
    }

    defer containerLogsResponse.Close()

    // Создание локального файла для записи логов
    file, err := os.Create("./local_directory/result")
    if err != nil {
        log.Fatalf("Failed to create file: %s", err.Error())
    }
    defer file.Close()

    // Запись логов в файл
    _, err = io.Copy(file, containerLogsResponse)
    if err != nil {
        log.Fatalf("Failed to write container logs to file: %s", err.Error())
    }

    // Удаление контейнера
    removeOptions := types.ContainerRemoveOptions{
        RemoveVolumes: true,
        Force:         true,
    }
    err = dockerClient.ContainerRemove(ctx, containerResponse.ID, removeOptions)
    if err != nil {
        log.Fatalf("Failed to remove container: %s", err.Error())
    }

    // Отправка файла на go сервер
    url := "http://localhost:4999/"
    filePath := "./local_directory/result"

    fileContents, err := ioutil.ReadFile(filePath)
    if err != nil {
        log.Fatalf("Failed to read file: %s", err.Error())
    }

    // Отправка POST запроса на сервер
    requestBody := &bytes.Buffer{}
    writer := multipart.NewWriter(requestBody)

    part, err := writer.CreateFormFile("file", filepath.Base(filePath))
    if err != nil {
        log.Fatalf("Failed to create form file: %s", err.Error())
    }

    _, err = io.Copy(part, bytes.NewReader(fileContents))
    if err != nil {
        log.Fatalf("Failed to copy file contents to form file: %s", err.Error())
    }

    contentType := writer.FormDataContentType()
    err = writer.Close()
    if err != nil {
        log.Fatalf("Failed to close multipart writer: %s", err.Error())
    }

    response, err := http.Post(url, contentType, requestBody)
    if err != nil {
        log.Fatalf("Failed to send POST request: %s", err.Error())
    }

    // Вывод статуса ответа от сервера
    log.Printf("Server response status: %s", response.Status)
