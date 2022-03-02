.PHONY: proto

proto:
	protoc -I ./proto \
		--go_out ./proto --go_opt paths=source_relative \
		--go-grpc_out ./proto --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative --grpc-gateway_opt logtostderr=true \
		./proto/captainhook/captainhook.proto

openapi:
	protoc -I ./proto \
		--openapiv2_out ./gen/openapiv2 \
		--openapiv2_opt logtostderr=true \
		./proto/captainhook/captainhook.proto

docker-up:
	docker-compose -f build/docker/docker-compose.yml up -d