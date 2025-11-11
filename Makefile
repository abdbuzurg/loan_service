# --- minimal config -----------------------------------------------------------
PROTO            ?= internal/proto/loan_service.proto
MIGRATIONS_DIR   ?= internal/platform/database/migrations
DB_DSN           ?= postgres://postgres:password@localhost:5432/asr_leasing?sslmode=disable

PROTOC           ?= protoc
SQLC             ?= sqlc
MIGRATE          ?= migrate

# --- proto generation ---------------------------------------------------------
.PHONY: proto
proto:
	$(PROTOC) -I . \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO)

# --- sqlc generation ----------------------------------------------------------
.PHONY: sqlc
sqlc:
	$(SQLC) generate

# --- migrations (golang-migrate) ----------------------------------------------
# create: make migrate-create name=add_payments_table
.PHONY: migrate-create
migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate-create name=<migration_name>"; exit 1)
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

# apply all up migrations
.PHONY: migrate-up
migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" up

# roll back one migration (change the step count if needed)
.PHONY: migrate-down
migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down 1

