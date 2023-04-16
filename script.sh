#bin/bash
docker run --name my-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=123 -e PGDATA=/pgdata -v /home/$USER/Documents/tinkoff/postgres/pgdata:/pgdata -p 5432:5432 -d postgres

cd ./docker_module/

docker build -f Dockerfile-pychrome -t pychrome .

python backend.py

curl -OL https://golang.org/dl/go1.16.7.linux-amd64.tar.gz

sudo tar -C /usr/local -xvf go1.16.7.linux-amd64.tar.gz

echo 'export PATH=$PATH:/usr/local/go/bin' >>  ~/.profile
cd ./Golang/
go run main.go

