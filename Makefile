all: build

build:
	@echo "Building..."
	@templ generate
	@echo "Generated templ templates..."
	@npx @tailwindcss/cli -i static/input.css -o static/styles.css
	@echo "Generated Tailwind CSS styles..."
	@go build -o temp/main cmd/api/main.go
	@echo "Generated Go binary..."

run:
	@go run cmd/api/main.go

clean:
	@echo "Cleaning..."
	@rm -rf temp/*

migrate-up:
	@echo "Migrating up..."
	@migrate -database "sqlite3://db.sqlite3" -path db/migration -verbose up

migrate-down:
	@echo "Migrating down..."
	@migrate -database "sqlite3://db.sqlite3" -path db/migration -verbose down
