from flask import Flask, request, jsonify
import docker

app = Flask(__name__)

@app.route('/run-tests', methods=['POST'])
def run_tests():
    tests = request.json['tests']
    for test in tests:
        # отправить задание в очередь для обработки
        # возвращаем пользователю ответ об успешном добавлении задания в очередь
    return jsonify({'message': 'Tests submitted successfully'})

@app.route('/get-results/<int:job_id>', methods=['GET'])
def get_results(job_id):
    # проверить, есть ли результаты выполнения этой задачи
    # если да, вернуть результаты, иначе сообщить, что результаты еще не готовы


@app.route('/jobs', methods=['GET'])
def get_jobs():
    # получить список всех задач из базы данных
    # вернуть список задач в формате JSON

@app.route('/jobs/<int:job_id>', methods=['GET'])
def get_job(job_id):
    # получить информацию о задаче из базы данных
    # вернуть информацию о задаче в формате JSON