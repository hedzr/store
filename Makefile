
test: store-test tests-test codecs-tests providers-tests

store-test:
	go test -v ./...

tests-test:
	cd ./tests && go test -v ./...

codecs-tests: hcl-test hjson-test json-test nestext-test toml-test yaml-test

hcl-test:
	cd ./codecs/hcl && go test -v ./...

hjson-test:
	cd ./codecs/hjson && go test -v ./...

json-test:
	cd ./codecs/json && go test -v ./...

nestext-test:
	cd ./codecs/nestext && go test -v ./...

toml-test:
	cd ./codecs/toml && go test -v ./...

yaml-test:
	cd ./codecs/yaml && go test -v ./...

providers-tests: env-test file-test

env-test:
	cd ./providers/env && go test -v ./...

file-test:
	cd ./providers/file && go test -v ./...







