fmt:
	gofmt -l -s -w .

run:
	go run main.go token_config.go

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o renovator .
