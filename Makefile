build:
	@mkdir -p ./bin
	@go build -o ./bin/vaal ./cmd/vaal/

install:
	@go install ./cmd/vaal/
