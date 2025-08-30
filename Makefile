all: build

build:
	@echo "Building..."
	@templ generate
	@tailwindcss -i static/input.css -o static/styles.css
	@go build -o temp/main cmd/api/main.go

run:
	@go run cmd/api/main.go

clean:
	@echo "Cleaning..."
	@rm -rf temp/*