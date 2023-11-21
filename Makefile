test:
	go test -v ./...

build:
	docker build -f ./docker/Dockerfile -t flights-server .

run:
	docker compose -f docker/docker-compose.yml up
