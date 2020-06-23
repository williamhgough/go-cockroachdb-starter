.env:
	cp .env.example .env

install-generators:
	go install \
		github.com/tmthrgd/go-bindata/go-bindata \
		github.com/vektra/mockery/cmd/mockery \

generate: install-generators
	go generate -x ./...

db-up:
	docker-compose up -d

db-down:
	docker-compose down

run: .env db-up
	go run main.go

test:
	go test ./... -cover -v -race
build:
	go build