swag:
	swag init -g cmd/main.go -o docs --parseDependency --parseInternal
up:
	docker-compose up -d
build:
	docker-compose build
down:
	docker-compose down