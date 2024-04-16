docker stop oracle-fil

docker rm oracle-fil

docker rmi oracle-fil

go mod tidy

docker build -t oracle-fil .

docker run --name oracle-fil -d oracle-fil

docker logs oracle-fil -f
