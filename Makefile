.PHONY: run
run:
	docker compose up

.PHONY: stop
stop:
	docker compose down

.PHONY: mysql-local
mysql-local:
	docker compose exec mysql mysql -u user -p go_layered_architecture

.PHONY: generate
generate:
	go generate ./...
