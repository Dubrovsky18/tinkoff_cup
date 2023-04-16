import shutil
import time
import docker
import os
import tempfile
from flask import Flask, request


app = Flask(__name__)

@app.route('/run-tests', methods=['POST'])
# Функция для запуска контейнера
def run_container():
    try:
        # Аргументы
        filename = request.args.get('filename')
        temp_dir = request.args.get('tempdir')

        # # Копируем файл в локальную папку докера
        # shutil.copy("/home/alla/Downloads/Telegram Desktop/webtest.py", '.')

        # Конфигурация Docker
        client = docker.from_env()

        # Сохраняем в переменной путь до файла, temp_dir - папка в которую бэк загрузил файл
        test_file_path = os.path.join(temp_dir, filename)

        # Создание Dockerfile из шаблона
        dockerfile_template_path = 'Dockerfile-template'
        path_to_test_in_container = f"/app/{filename}"
        # print(test_file_path)
        # if os.path.isfile(test_file_path):
        #     print('File exist')
        # else:
        #     print('File dont exist')

        # Добавление в шаблон кастомных строк
        dockerfile_copy_path = os.path.join(temp_dir, 'Dockerfile')
        with open(dockerfile_template_path, 'r') as f:
            dockerfile_content = f.read()
            dockerfile_content += f"\nCOPY {filename} /app/{filename}\n"
            dockerfile_content += f"\nCMD [\"python\", \"{path_to_test_in_container}\"]\n "

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

        # Удаление временной директории с файлами
        # shutil.rmtree(test_file_path, False)



    except Exception as e:
        print('Error:', e)


# run_container()
