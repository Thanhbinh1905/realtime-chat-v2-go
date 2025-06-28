include .env

DB_URL=postgres://binh:secret@localhost:5432/auth_db?sslmode=disable


migrate-up:
	migrate -path ./auth-service/internal/db/migrations -database "$(DB_URL)" up
migrate-down:
	migrate -path ./auth-service/internal/db/migrations -database "$(DB_URL)" down