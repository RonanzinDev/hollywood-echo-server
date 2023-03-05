build:
	go build -o bin/actor


run: build
	./bin/actor