BUILD = dist
PACKAGE_NAME = github.com/iortego42/go-rat

build_client: build_dir
	go build -o $(BUILD)/client client/client.go

build_implant: build_dir
	go build -o $(BUILD)/implant implant/implant.go

build_server: build_dir
	go build -o $(BUILD)/server server/server.go

build_dir:
	@[[ -d $(BUILD) ]] || mkdir $(BUILD)

client: build_client
	$(BUILD)/client whoami

implant: build_implant
	$(BUILD)/implant whoami

server: build_server
	$(BUILD)/server whoami

protoc:
	protoc --go_out=. --go_opt=module=$(PACKAGE_NAME) --go-grpc_out=module=$(PACKAGE_NAME):.  grpcapi/implant.proto

.PHONY: all server client implant build_dir build_client build_implant build_server