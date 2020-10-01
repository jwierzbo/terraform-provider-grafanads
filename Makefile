NAME="terraform-provider-grafanads"
VERSION="1.0.0"
BINARY=${NAME}

.PHONY: terraform-provider-grafanads
terraform-provider-grafanads:
	go install github.com/jwierzbo/terraform-provider-grafanads/cmd/${NAME}

.PHONY: update-libs
update-libs:
	GOSUMDB=off GOPROXY=direct go mod tidy -v
	go mod vendor -v

test:
	go test ./...

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: create-release
create-release:
	mkdir -p build/${NAME}/linux-amd64
	mkdir -p build/${NAME}/darwin-amd64

	GOOS=linux GOARCH=amd64 go build -o build/${NAME}/linux-amd64/${NAME} \
		github.com/jwierzbo/terraform-provider-grafanads/cmd/${NAME}

	GOOS=darwin GOARCH=amd64 go build -o build/${NAME}/darwin-amd64/${NAME} \
		github.com/jwierzbo/terraform-provider-grafanads/cmd/${NAME}

.PHONY: install-local
install-local:
	go build -o ${BINARY} github.com/jwierzbo/terraform-provider-grafanads/cmd/${NAME}
	mkdir -p ~/.terraform.d/plugins/
	mv ${BINARY} ~/.terraform.d/plugins/
