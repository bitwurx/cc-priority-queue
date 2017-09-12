.PHONY: build
build:
	@docker run \
		--rm \
		-e CGO_ENALBED=0 \
		-e GOOS=linux \
		-v $(PWD):/usr/src/concord-pq \
		-w /usr/src/concord-pq \
		golang go get -v -d && go build -o main
	@docker build -t concord/pq .

.PHONY: test
test:
	@docker run \
		-d \
		-e ARANGO_ROOT_PASSWORD=abc123 \
		--name concord-pq_test__arangodb \
		arangodb/arangodb
	@docker run \
		-d \
		-e ARANGODB_HOST=http://arangodb:8529 \
		-e ARANGODB_NAME=test__concord_pq \
		-e ARANGODB_USER=root \
		-e ARANGODB_PASS=abc123 \
		-v $(PWD):/go/src/concord-pq \
		-v $(PWD)/.src:/go/src \
		-w /go/src/concord-pq \
		--link concord-pq_test__arangodb:arangodb \
		--name concord-pq_test \
		golang /bin/sh -c "go get -v -t -d ./... && go test -v ./..."
	@docker logs -f concord-pq_test
	@docker rm -f concord-pq_test
	@docker rm -f concord-pq_test__arangodb

.PHONY: test-short
test-short:
	@docker run \
		--rm \
		-v $(PWD):/go/src/concord-pq \
		-v $(PWD)/.src:/go/src \
		-w /go/src/concord-pq \
		golang /bin/sh -c \
			"mkdir -p .coverage && \
			go get -v -t -d ./... && \
			go test -short -v -cover -coverprofile=.coverage/queue.out ./queue && \
			go test -short -v -cover -coverprofile=.coverage/task.out ./task"
