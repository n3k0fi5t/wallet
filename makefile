all: build
build: clean
	docker-compose build

clean:
	docker-compose down
	rm -rf ./data

run: build
	docker-compose up

.PHONY: all build clean run