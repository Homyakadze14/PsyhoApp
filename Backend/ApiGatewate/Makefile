gen-proto:
	protoc -I proto proto/${pkg}/*.proto --go_out=proto/gen/ --go_opt=paths=source_relative --go-grpc_out=proto/gen/ --go-grpc_opt=paths=source_relative

gen-docs:
	swag init -g=internal/controller/rest/v1/router.go -o=docs --parseInternal

run:
	go run cmd/app/main.go --config=config/local.yaml
