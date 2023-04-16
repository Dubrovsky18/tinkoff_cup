import shutil
import time
import docker
import os
import tempfile

@app.route('/run-tests?тфьу', methods=['POST'])
# Функция для запуска контейнера
def run_container():
    try:
        # Конфигурация Docker
        client = docker.from_env()

        # 1. Запись временного файла с тестами
        temp_dir = tempfile.mkdtemp()
        test_file_path = os.path.join(temp_dir, 'test.py')
        with open(test_file_path, 'w') as f:
            f.write('print("Hello world")')

        # 2. Создание Dockerfile-template с копированием файла с тестами
        dockerfile_template_path = 'Dockerfile-template'
        path_to_test_in_container = '/app/test.py'
        print(test_file_path)
        if os.path.isfile(test_file_path):
            print('File exist')
        else:
            print('File dont exist')

        shutil.copy(test_file_path, '.')

        dockerfile_copy_path = os.path.join(temp_dir, 'Dockerfile')
        with open(dockerfile_template_path, 'r') as f:
            dockerfile_content = f.read()
            dockerfile_content += f"\nCOPY test.py /app/test.py\n" \
                                  f"\nCMD [\"python\", \"{path_to_test_in_container}\"]\n"
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

        # 3. Удаление временной директории с файлами
        # temp_dir.cleanup()

        # time.sleep(600)
        # container.stop()
        # container.remove()



    except Exception as e:
        print('Error:', e)


run_container()
