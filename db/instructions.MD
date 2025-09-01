1. Install `golang-migrate` using `go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest` with sqlite3 driver
2. Run `migrate create -ext sql -dir db/migration -seq init_schema` to create migrations
3. Fill in the sql files
4. Run `migrate -database "sqlite3://db.sqlite3" -path db/migration -verbose up` to migrate
