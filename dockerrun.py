import docker
import pika
import json
import psycopg2
import threading

# Конфигурация Docker
client = docker.from_env()

# Конфигурация RabbitMQ
connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost'))
channel = connection.channel()
channel.queue_declare(queue='test_queue')

# Конфигурация PostgreSQL
conn = psycopg2.connect(database="mydatabase", user="myusername", password="mypassword", host="localhost", port="5432")
cur = conn.cursor()

# Функция для запуска контейнера
def run_container(test_id, dockerfile_path):
    try:
        # Записываем информацию о тесте в базу данных
        cur.execute("INSERT INTO tests (id, status) VALUES (%s, 'running')", (test_id,))
        conn.commit()

        # Создаем контейнер и запускаем тест
        container = client.containers.run(image=dockerfile_path, command="python /app/run_tests.py {}".format(test_id), detach=True)

        # Отправляем сообщение с информацией о контейнере в очередь
        message = {'test_id': test_id, 'container_id': container.id}
        channel.basic_publish(exchange='', routing_key='test_queue', body=json.dumps(message))

    except Exception as e:
        print('Error:', e)

# Функция для чтения сообщений из очереди
def consume_queue():
    def callback(ch, method, properties, body):
        message = json.loads(body)

        # Записываем информацию о контейнере в базу данных
        cur.execute("UPDATE tests SET status = 'completed', result = %s, logs = %s WHERE id = %s",
                    (message['result'], message['logs'], message['test_id']))
        conn.commit()

        # Остановка и удаление контейнера
        container = client.containers.get(message['container_id'])
        container.stop()
        container.remove()

    channel.basic_consume(queue='test_queue', on_message_callback=callback, auto_ack=True)
    channel.start_consuming()

# Функция для получения тестов из веб-интерфейса и добавления их в очередь
def add_tests_to_queue(tests):
    for test in tests:
        # Записываем информацию о тесте в базу данных
        cur.execute("INSERT INTO tests (name, status) VALUES (%s, 'pending')", (test,))
        conn.commit()

        # Добавляем тест в очередь
        message = {'test_id': cur.lastrowid, 'test_name': test}
        channel.basic_publish(exchange='', routing_key='test_queue', body=json.dumps(message))

# Запускаем чтение очереди в отдельном потоке
threading.Thread(target=consume_queue).start()