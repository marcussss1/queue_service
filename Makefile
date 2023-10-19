.PHONY: build
build:
	docker build --tag api-image .

.PHONY: run
run:
	docker run --name=api -p 8080:8080 api-image

.PHONY: stop
stop: |
	docker stop api
	docker rm api
