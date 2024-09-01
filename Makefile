BUILD = dist
PACKAGE_NAME = github.com/iortego42/go-rat

build_client: $(BUILD)
	go build -o $(BUILD)/client ./client

build_implant: $(BUILD)
	go build -o $(BUILD)/implant ./implant

build_server: $(BUILD)
	go build -o $(BUILD)/server ./server

$(BUILD):
	mkdir $(BUILD)

client: build_client
	$(BUILD)/client whoami

implant: build_implant
	$(BUILD)/implant

server: build_server
	$(BUILD)/server

protoc:
	protoc --go_out=. --go_opt=module=$(PACKAGE_NAME) --go-grpc_out=module=$(PACKAGE_NAME):.  grpcapi/implant.proto

.PHONY: all server client implant build_dir build_client build_implant build_server
