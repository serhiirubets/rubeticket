# Best Concert Ticket App

### How to run project:
- git clone https://github.com/serhiirubets/rubeticket.git
- Run docker-compose file: `docker compose up -d`
- Install all dependencies `go mod tidy`
- Run migration once `go run migrations/auto.go`
- Run app `go run cmd/main.go`

### Generate swagger: `make swagger`
It will generate swagger that you can see if open http://localhost:7777/swagger/index.html