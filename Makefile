.PHONY: generate
generate:
	rm -rf ./app/mock/*
	go generate ./...
