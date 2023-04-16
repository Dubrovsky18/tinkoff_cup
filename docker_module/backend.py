import shutil
import time
import docker
import os
import tempfile

import requests
from flask import Flask, request

app = Flask(__name__)

@app.route('/run-tests', methods=['POST'])
# Функция для запуска контейнера
def run_container():
    try:
        # Аргументы
        filename = request.args.get('filename')
        temp_dir = request.args.get('tempdir')

        # Конфигурация Docker
        client = docker.from_env()

        # Сохраняем в переменной путь до файла, temp_dir - папка в которую бэк загрузил файл
        test_file_path = os.path.join(temp_dir, filename)

        # Создание Dockerfile из шаблона
        dockerfile_template_path = 'Dockerfile-template'
        path_to_test_in_container = f"/app/{filename}"

        # Добавление в шаблон кастомных строк
        dockerfile_copy_path = os.path.join(temp_dir, 'Dockerfile')
        with open(dockerfile_template_path, 'r') as f:
            dockerfile_content = f.read()
            dockerfile_content += f"\nCOPY {filename} /app/{filename}\n"
            dockerfile_content += f"\nCMD [\"python\", \"{path_to_test_in_container}\", \" >> result\"]\n "

        with open(dockerfile_copy_path, 'w') as f:
            f.write(dockerfile_content)

        print(dockerfile_content)

        # Создаем контейнер и запускаем тест
        image, build_logs = client.images.build(path=temp_dir)
        for line in build_logs:
            print(line)

        container = client.containers.run(image.id, detach=True, stream=True)
        for log in container.logs(stream=True):
            print(log)

        # Копируем из контейнера файл с логами
        os.makedirs("./local_directory", exist_ok=True)

        with open("./local_directory/result", "wb") as f:
            bits, _ = container.get_archive("/app/result")
            for chunk in bits:
                f.write(chunk)

        # Удаляем контейнер
        container.remove(force=True)

        # Отправляем файл на go сервер
        url = "http://localhost:4999/"
        file_path = "./local_directory/result"

        with open(file_path, "rb") as file:
            response = requests.post(url, files={"file": file})
        print(response.status_code)

    except Exception as e:
        print('Error:', e)


if __name__ == '__main__':
    app.run(port=5000)
