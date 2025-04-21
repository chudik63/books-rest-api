swag:
	swag init -g cmd/main.go -o docs --parseDependency --parseInternal
postgres-up:
	docker-compose up -d
run:
	go run cmd/main.go