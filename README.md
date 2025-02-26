### First project Oauth with go lang

 - access_token expires in 1 minuto (to facility test of the refresh)
 - refresh_token expires in 7 days<br/>
 You can change this in `core/helpers/auth/create_token.go`

### Dependencies

go get github.com/jackc/pgx
go get github.com/jackc/pgx/v5/pgconn@v5.7.2

go get github.com/lib/pq

go get github.com/joho/godotenv <br/>
```env
DB_HOST=127.0.0.1 or localhost  ...
DB_USER=<db-user>
DB_PASS=<db-pass>
DB_NAME=<db-name>
DB_PORT=<db-port> postgres default 5432
API_PORT=<api-port> 8000 3000 ...
DB_SSL=disable
DB_TZ=America/Sao_Paulo your timezone
SECRET_KEY=<secre-key>
REFRESH_SECRET_KEY=<refresh-secret-key>
```

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
