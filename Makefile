build:
	CGO_ENABLED=0 go build -o ./example ./main.go

test:
	go test -v ./...

run:
	go run main.go

mock:
	mockery --all