import (
    "fmt"
    "net/http"
    "os"
    "io/ioutil"

    "github.com/docker/docker/client"
)

func runContainer(w http.ResponseWriter, r *http.Request) {
    // Аргументы
    filename := r.FormValue("filename")
    temp_dir := r.FormValue("tempdir")

    // Конфигурация Docker
    cli, err := client.NewEnvClient()
    if err != nil {
        fmt.Println(err)
        return
    }

    // Сохраняем в переменной путь до файла, temp_dir - папка в которую бэк загрузил файл
    test_file_path := filepath.Join(temp_dir, filename)

    // Создание Dockerfile из шаблона
    dockerfile_template_path := "Dockerfile-template"
    path_to_test_in_container := fmt.Sprintf("/app/%s", filename)

    // Добавление в шаблон кастомных строк
    dockerfile_content, err := ioutil.ReadFile(dockerfile_template_path)
    if err != nil {
        fmt.Println(err)
        return
    }

    dockerfile_content += fmt.Sprintf("\nCOPY %s /app/%s\n", filename, filename)
    dockerfile_content += fmt.Sprintf("\nCMD [\"python\", \"%s\", \" >> result\"]\n", path_to_test_in_container)

    dockerfile_copy_path := filepath.Join(temp_dir, "Dockerfile")
    err = ioutil.WriteFile(dockerfile_copy_path, []byte(dockerfile_content), 0644)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(dockerfile_content)

    // Создаем контейнер и запускаем тест
    ctx := context.Background()
    options := types.ImageBuildOptions{
        Context:    ctx,
        Dockerfile: "Dockerfile",
        Remove:     true,
        ForceRemove:true,
        NoCache:    true,
        PullParent: true,
        Tags:       []string{"test-image:latest"},
    }

    buildResponse, err := cli.ImageBuild(ctx, buildContext, options)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer buildResponse.Body.Close()
    io.Copy(os.Stdout, buildResponse.Body)

    container, err := cli.ContainerCreate(ctx, &container.Config{
        Image: "test-image:latest",
    }, nil, nil, "")
    if err != nil {
        fmt.Println(err)
        return
    }

    err = cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
    if err != nil {
        fmt.Println(err)
        return
    }

    statusCh, errCh := cli.ContainerWait(ctx, container.ID, container.WaitConditionNextExit)

    select {
    case err := <-errCh:
        if err != nil {
            fmt.Println(err)
            return
        }
    case <-statusCh:
    }

    out, err := cli.ContainerLogs(ctx, container.ID, types.ContainerLogsOptions{ShowStdout: true})
    if err != nil {
        fmt.Println(err)
        return
    }

    stdcopy.StdCopy(os.Stdout, os.Stderr, out)

    // Копируем из контейнера файл с логами
    os.MkdirAll("./local_directory", 0755)

    bits, _, err
