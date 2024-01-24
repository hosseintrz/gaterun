migrate_up:
	migrate -path pkg/database/postgres/migrations/ -database "postgresql://postgres:pass1234@localhost:5432/gaterun?sslmode=disable" -verbose up

migrate_down:
	migrate -path pkg/database/postgres/migrations/ -database "postgresql://postgres:pass1234@localhost:5432/gaterun?sslmode=disable" -verbose down

migrate_force:
	migrate -path pkg/database/postgres/migrations/ -database "postgresql://postgres:pass1234@localhost:5432/gaterun?sslmode=disable" force $(version)

create_migration:
	migrate create -ext sql -dir pkg/database/migrations -seq $(name)

.PHONY: migrate_up migrate_down migrate_force create_migration 