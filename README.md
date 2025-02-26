### Dependencies

go get github.com/jackc/pgx
go get github.com/jackc/pgx/v5/pgconn@v5.7.2

go get github.com/lib/pq

go get github.com/joho/godotenv

go get github.com/gorilla/mux

go get github.com/google/uuid

go get github.com/badoux/checkmail

go get golang.org/x/crypto/bcrypt

go get github.com/dgrijalva/jwt-go

go get github.com/golang-migrate/migrate

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

### Create migration up and down
migrate create -ext sql -dir ./core/data/db/migrations <table_name>

### Migrate up to db
migrate -database postgres://DB_USER:DB_PASS@DB_HOST:DB_PORT/DB_NAME?sslmode=disable -path ./core/data/db/migrations up

### Migrate down db
migrate -database postgres://DB_USER:DB_PASS@DB_HOST:DB_PORT/DB_NAME?sslmode=disable -path ./core/data/db/migrations down

### Run
`go run main.go`
