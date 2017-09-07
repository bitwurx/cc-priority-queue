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
		--rm \
		-v $(PWD):/usr/src/concord-pq \
		-v $(PWD)/.src:/go \
		-w /usr/src/concord-pq \
		golang \
			mkdir -p .coverage && \
			go get -v -t -d && \
			go test -v -cover -coverprofile=.coverage/queue.out ./queue && \
			go test -v -cover -coverprofile=.coverage/task.out  ./task