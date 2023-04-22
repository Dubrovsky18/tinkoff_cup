package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "path/filepath"
    "strings"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/api/types/mount"
    "github.com/docker/docker/client"
    "github.com/docker/docker/pkg/archive"
    "github.com/docker/docker/pkg/stdcopy"
    "github.com/gorilla/mux"
)

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/run-tests", runTests).Methods("POST")
    http.ListenAndServe(":5000", router)
}

func runTests(w http.ResponseWriter, r *http.Request) {
    // Аргументы
    filename := r.URL.Query().Get("filename")
    tempDir := r.URL.Query().Get("tempdir")

    // Конфигурация Docker
    cli, err := client.NewEnvClient()
    if err != nil {
        fmt.Println("Error creating Docker client:", err)
        return
    }

    // Сохраняем в переменной путь до файла, tempDir - папка в которую бэк загрузил файл
    testFilePath := filepath.Join(tempDir, filename)

    // Создание Dockerfile из шаблона
    dockerfileTemplatePath := "Dockerfile-template"
    pathToTestInContainer := fmt.Sprintf("/app/%s", filename)

    // Добавление в шаблон кастомных строк
    dockerfileCopyPath := filepath.Join(tempDir, "Dockerfile")
    dockerfileContent, err := ioutil.ReadFile(dockerfileTemplatePath)
    if err != nil {
        fmt.Println("Error reading Dockerfile template:", err)
        return
    }
    dockerfileContent = append(dockerfileContent, []byte(fmt.Sprintf("\nCOPY %s %s\n", filename, pathToTestInContainer))...)
    dockerfileContent = append(dockerfileContent, []byte(fmt.Sprintf("\nCMD [\"python\", \"%s\", \" >> result\"]\n ", pathToTestInContainer))...)
    if err = ioutil.WriteFile(dockerfileCopyPath, dockerfileContent, 0666); err != nil {
        fmt.Println("Error writing Dockerfile:", err)
        return
    }

    fmt.Println(string(dockerfileContent))

    // Создаем контейнер и запускаем тест
    imageBuildResponse, err := cli.ImageBuild(r.Context(), archive.ReadDir(tempDir), types.ImageBuildOptions{
        Remove: true,
        Tags:   []string{"my-tag"},
    })
    if err != nil {
        fmt.Println("Error building Docker image:", err)
        return
    }
    defer imageBuildResponse.Body.Close()

    _, err = stdcopy.StdCopy(os.Stdout, os.Stderr, imageBuildResponse.Body)
    if err != nil {
        fmt.Println("Error copying Docker image build response:", err)
        return
    }

    resp, err := cli.ContainerCreate(r.Context(), &container.Config{
        Image: "my-tag",
        Tty:   true,
        Cmd:   []string{"python", pathToTestInContainer},
    }, &container.HostConfig{
        Mounts: []mount.Mount{
            {
                Type:   mount.TypeBind,
                Source: tempDir,
                Target: "/app",
            },
        },
    }, nil, "")
    if err != nil {
        fmt.Println("Error creating Docker container:", err)
        return
    }

