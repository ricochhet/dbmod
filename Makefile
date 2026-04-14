BUILD_OUTPUT=build
ASSET_PATH=assets

CUSTOM=-X 'main.buildDate=$(shell date)' -X 'main.gitHash=$(shell git rev-parse --short HEAD)' -X 'main.buildOn=$(shell go version)'
LDFLAGS=$(CUSTOM) -w -s -extldflags=-static
GO_BUILD=go build -trimpath -ldflags "$(LDFLAGS)"

APP_NAMES=dbmod

APP_PATH=.
APP_BIN_NAME=dbmod

define GO_BUILD_APP
	CGO_ENABLED=1 GOOS=$(1) GOARCH=$(2) $(GO_BUILD) -o $(BUILD_OUTPUT)/$(3) $(4)
endef

.PHONY: all
all: dbmod

.PHONY: fmt
fmt:
	gofumpt -l -w -extra .

.PHONY: tidy
tidy:
	@echo "[main] tidy"
	go mod tidy

.PHONY: update
update:
	@echo "[main] update dependencies"
	go get -u ./...

.PHONY: lint
lint: fmt
	@echo "[main] golangci-lint"
	golangci-lint run ./... --fix

.PHONY: test
test:
	go test ./...

.PHONY: deadcode
deadcode:
	deadcode ./...

.PHONY: syso
syso:
	windres $(APP_PATH)/app.rc -O coff -o $(APP_PATH)/app.syso

.PHONY: png-to-icos
png-to-icos:
	magick $(ASSET_PATH)/win-icon.png -background none -define icon:auto-resize=256,128,64,48,32,16 $(ASSET_PATH)/win-icon.ico

.PHONY: copy-assets
copy-assets:
	cp -r $(ASSET_PATH)/* $(BUILD_OUTPUT)

.PHONY: gen-certs
gen-certs:
	mkcert localhost 127.0.0.1 ::1

# ----- dbmod -----
.PHONY: dbmod
dbmod: dbmod-linux dbmod-linux-arm64 dbmod-darwin dbmod-darwin-arm64 dbmod-windows

.PHONY: dbmod-linux
dbmod-linux: fmt
	$(call GO_BUILD_APP,linux,amd64,$(APP_BIN_NAME)-linux,$(APP_PATH))

.PHONY: dbmod-linux-arm64
dbmod-linux-arm64: fmt
	$(call GO_BUILD_APP,linux,arm64,$(APP_BIN_NAME)-linux-arm64,$(APP_PATH))

.PHONY: dbmod-darwin
dbmod-darwin: fmt
	$(call GO_BUILD_APP,darwin,amd64,$(APP_BIN_NAME)-darwin,$(APP_PATH))

.PHONY: dbmod-darwin-arm64
dbmod-darwin-arm64: fmt
	$(call GO_BUILD_APP,darwin,arm64,$(APP_BIN_NAME)-darwin-arm64,$(APP_PATH))

.PHONY: dbmod-windows
dbmod-windows: fmt copy-assets
	$(call GO_BUILD_APP,windows,amd64,$(APP_BIN_NAME).exe,$(APP_PATH))
