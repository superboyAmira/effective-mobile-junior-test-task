
# EXTERNAL API
EXTERNAL_API_DIR=./external_api
EXTERNAL_API_OPENAPI=./api/external_openapi.yml

codegen:
	# sudo npm install @openapitools/openapi-generator-cli -g
	openapi-generator-cli generate -i $(EXTERNAL_API_OPENAPI) -g go -o $(EXTERNAL_API_DIR)

container_up :
	docker-compose -f docker-compose.yaml up

container_rm:
	docker-compose stop \
	&& docker-compose rm
