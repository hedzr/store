
# TEST_OPTS = -v -cover -race -coverprofile=./logs/coverage$(seq).txt -covermode=atomic -timeout=20m -test.short -vet=off -test.v
TEST_OPTS = -v -cover -race -covermode=atomic -timeout=20m -vet=off -test.v




test: store-test tests-test codecs-tests providers-tests

# go test -v -cover -test.v -test.run '^$(tname)$$' $(package) -race -coverprofile=./logs/coverage$(seq).txt -covermode=atomic -timeout=20m -test.short -vet=off

store-test:
	go test $(TEST_OPTS) -coverprofile=./logs/coverage-$(seq).txt ./...

tests-test:
	cd ./tests && go test $(TEST_OPTS)  ./...

codecs-tests: hcl-test hjson-test json-test nestext-test toml-test yaml-test

hcl-test:
	cd ./codecs/hcl && go test $(TEST_OPTS) ./...

hjson-test:
	cd ./codecs/hjson && go test $(TEST_OPTS) ./...

json-test:
	cd ./codecs/json && go test $(TEST_OPTS) ./...

nestext-test:
	cd ./codecs/nestext && go test $(TEST_OPTS) ./...

toml-test:
	cd ./codecs/toml && go test $(TEST_OPTS) ./...

yaml-test:
	cd ./codecs/yaml && go test $(TEST_OPTS) ./...

providers-tests: env-test file-test

env-test:
	cd ./providers/env && go test $(TEST_OPTS) ./...

file-test:
	cd ./providers/file && go test $(TEST_OPTS) ./...







