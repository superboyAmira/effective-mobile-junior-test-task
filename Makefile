
# EXTERNAL API

EXTERNAL_API_DIR=./external_api
EXTERNAL_API_OPENAPI=./api/external_openapi.yml

codegen:
	# sudo npm install @openapitools/openapi-generator-cli -g
	sopenapi-generator-cli generate -i openapi.yml -g go-server -o .

container_up :
	docker-compose -f docker-compose.yaml up

container_rm:
	docker-compose stop \
	&& docker-compose rm
