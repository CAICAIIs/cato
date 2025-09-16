

build-proto:
	protoc -I=./proto --go_out=../ ./proto/*.proto


install:
	go install ./cmd/protoc-gen-cato