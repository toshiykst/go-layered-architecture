.PHONY: generate
generate:
	rm -rf ./app/mock/*
	go generate ./...

.PHONY: mysql-local
mysql-local:
	docker compose exec mysql mysql -u user -p go_layered_architecture
