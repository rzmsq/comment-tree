.PHONY: quality lint vet

quality: lint vet

lint:
	@echo "Running linter..."
	golangci-lint run ./...

vet:
	@echo "Running vet..."
	go vet ./...

run:
	docker-compose down -v
	docker-compose build --no-cache
	docker-compose up

tools:
	go get github.com/ilyakaznacheev/cleanenv
	go get github.com/go-playground/validator
	go get github.com/mattn/go-sqlite3
	go get github.com/wb-go/wbf@v0.0.8