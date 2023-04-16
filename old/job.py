from flask import Flask, request, jsonify
import docker

app = Flask(__name__)
client = docker.from_env()

@app.route('/run_tests', methods=['POST'])
def run_tests():
    image_name = request.json['image_name']
    command = request.json['command']
    container = client.containers.run(image_name, command)
    return jsonify({'container_id': container.id})

@app.route('/container_logs/<container_id>', methods=['GET'])
def container_logs(container_id):
    container = client.containers.get(container_id)
    return container.logs()

if __name__ == '__main__':
    app.run(debug=True)