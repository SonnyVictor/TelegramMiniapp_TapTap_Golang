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


proto:
	rm -f pb/*.go
	protoc --proto_path=internal/proto --go_out=internal/pb --go_opt=paths=source_relative \
	--go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/pb --grpc-gateway_opt=paths=source_relative \
	internal/proto/*.proto


evans:
	evans --host localhost --port 9090 -r repl

.PHONY: server migrateup new_migration migratedown migratedown1 proto evans