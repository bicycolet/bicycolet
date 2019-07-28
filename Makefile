REPO_PATH = github.com/bicycolet/bicycolet

setup:
	go get -u $(REPO_PATH)...

build: bin/bicycolet

install:
	go install -v $(REPO_PATH)/cmd/bicycolet/...

clean:
	@rm -f bin/bicycolet 2> /dev/null || true

bin/bicycolet:
	go build -o bin/bicycolet $(REPO_PATH)/cmd/bicycolet

check: clean build install
	@cd test && ./main.sh
