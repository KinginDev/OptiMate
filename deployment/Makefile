USER_SERVICE_BINARY=./bin/userService

build_user_service:
	@echo "Building user service"
	cd ../user-service && env GOOS=linux CGO_ENABLED=0 go build -o ${USER_SERVICE_BINARY} ./cmd/api/main.go
	@echo "User service build completed"

build_backend: build_user_service
	@echo "Building backend"
	docker-compose up --build -d
	@echo "Backend build completed"

down:
	@echo "Stopping backend"
	docker-compose down
	@echo "Backend stopped"