.PHONY: migrate

FILE ?= 2026-01-15-17-00.sql

migrate:
	docker exec -i go-cs-postgres-1 psql -U postgres -d postgres < migrations/$(FILE)
