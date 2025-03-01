DB_URL=postgresql://root:secret@localhost:5432/simple_telegram?sslmode=disable



server:
	go run cmd/app/main.go


migrateup:
	migrate -path pkg/migration -database "postgresql://root:secret@localhost:5432/simple_telegram?sslmode=disable" up 

new_migration:
	migrate create -ext sql -dir pkg/migration -seq $(name)

migratedown:
	migrate -path pkg/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path pkg/migration -database "$(DB_URL)" -verbose down 1

.PHONY: server migrateup new_migration migratedown migratedown1