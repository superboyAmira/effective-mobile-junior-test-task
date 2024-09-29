
# EXTERNAL API
EXTERNAL_API_DIR=./external_api
EXTERNAL_API_OPENAPI=./api/external_openapi.yml

swag_docs:
	@export PATH=$PATH:$(go env GOPATH)/bin 
	@swag init --dir ./cmd/online-song-library,./internal/

external_api_gen:
	# sudo npm install @openapitools/openapi-generator-cli -g
	openapi-generator-cli generate -i $(EXTERNAL_API_OPENAPI) -g go-server -o $(EXTERNAL_API_DIR)

container_up :
	docker-compose -f docker-compose.yaml up

container_rm:
	docker-compose stop \
	&& docker-compose rm

test:
	go test -v ./test/*.go
	# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	# все линтеры в одном месте
	golangci-lint run