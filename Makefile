include .env
export

TOOLS_PATH=bin/tools
ogen=$(TOOLS_PATH)/ogen

$(ogen):
	GOBIN=`pwd`/$(TOOLS_PATH) go install github.com/ogen-go/ogen/cmd/ogen@v1.5.0

setup: $(ogen)

.PHONY: run
run:
	go run ./cmd/cdn

.PHONY: compose.up
compose.up:
	docker compose up -d

.PHONY: compose.down
compose.down:
	docker compose down

.PHONY: generate.api
generate.api: $(ogen)
	$(ogen) --loglevel error --clean --config .ogen.yml --target ./api/gen ./api/openapi.yaml
	docker run --rm -v `pwd`:/spec redocly/cli:1.25.3 bundle ./api/openapi.yaml > ./api/bundle.yaml
