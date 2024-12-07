USER_SERVICE_BINARY=./bin/userService
OPTIMIZER_SERVICE_BINARY=./bin/optimizerService

build_user_service:
	@echo "Generate swagger docs for user service"
	cd ../user-service && swag init --generalInfo ./cmd/api/main.go  -o ./cmd/api/docs
	@echo "Swagger docs generated for user service"
	@echo "Building user service"
	cd ../user-service && env GOOS=linux CGO_ENABLED=0 go build -o ${USER_SERVICE_BINARY} ./cmd/api/main.go
	@echo "User service build completed"

build_optimizer_service:
	@echo "Generate swagger docs for optimizer service"
	cd ../optimizer-service && swag init --generalInfo ./cmd/api/main.go  -o ./cmd/api/docs
	@echo "Swagger docs generated for optimizer service"
	@echo "Building optimizer service"
	cd ../optimizer-service && env GOOS=linux CGO_ENABLED=0 go build -o ${OPTIMIZER_SERVICE_BINARY} ./cmd/api/main.go
	@echo "Optimizer service build completed"

build: build_user_service build_optimizer_service
	@echo "Building backend"
	docker-compose up --build -d
	@echo "Backend build completed"

down:
	@echo "Stopping backend"
	docker-compose down
	@echo "Backend stopped"