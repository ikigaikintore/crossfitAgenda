.PHONY: tests

help: ## Show this help
	@echo "Help"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[93m %s\n", $$1, $$2}'

install: ## Install dependencies
	go mod tidy

generate-api-v1: generate-http-api-v1 generate-proto-v1 install ## Generate http API v1

generate-http-api-v1: ## Generates the endpoint for full v1 API
	rm -f ./service/handlers/http/server.gen.go && \
	mkdir -p ./out/api/ && \
	docker run \
	--rm -v ${PWD}:/local \
	docker.io/ervitis/oapi-codegen:v1.12.4 oapi-codegen --config /local/service/handlers/oapiconfig.yaml /local/service/handlers/schema.yml > ./out/api/server.gen.go && \
	sleep 3 && \
	mv ./out/api/server.gen.go ./service/handlers/http/ && \
	sleep 1 && rm -rf ./out

# The proto file generated:
# - has to be in one file
# - change the package name to "v1"
# - add the option "go_package = \"./;v1\"
generate-proto-v1: ## Generate proto file WARNING this is beta and you have to move content into one file, then generate the protobuf implementations with 'generate-grpc-v1'
	rm -rf ./service/handlers/grpc/*.proto && \
	docker run \
	--rm -v ${PWD}:/local \
	docker.io/openapitools/openapi-generator-cli:v6.2.1 generate \
	-g protobuf-schema \
	-o /local/service/handlers/grpc/ \
	--package-name crossfitagenda \
	-i /local/service/handlers/schema.yml && \
	touch ./service/handlers/grpc/crossfitagenda.proto && \
	cat ./service/handlers/grpc/models/{agenda,error,status}.proto >> ./service/handlers/grpc/crossfitagenda.proto && \
	cat ./service/handlers/grpc/services/default_service.proto >> ./service/handlers/grpc/crossfitagenda.proto && \
	rm -rf ./service/handlers/grpc/models ./service/handlers/grpc/services ./service/handlers/grpc/README.md ./service/handlers/grpc/.openapi-generator-ignore ./service/handlers/grpc/.openapi-generator

generate-grpc-v1: ## Generate the server and client from proto file
	rm -rf ./service/handlers/grpc/*.pb.go && \
	docker run \
	--rm -v ${PWD}:/defs namely/protoc-all \
	-f ./service/handlers/grpc/*.proto \
	-l go \
	--go_out=plugins=grpc:/gen/grpc && \
	sleep 1 && \
	mv ./gen/pb-go/* ./service/handlers/grpc/ && \
	sleep 1 && \
	rm -rf ./gen