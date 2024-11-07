# Makefile tested working on Ubuntu Server 24.04 LTS
BINARY_NAME=server-shutdown

build:
	go build -o ./bin/${BINARY_NAME}

run: build
	./bin/${BINARY_NAME}

clean:
	go clean
	rm -f ./bin

install: build clean
	cp ./bin/${BINARY_NAME} /usr/local/bin
	chown root: /usr/local/bin/${BINARY_NAME}
	chmod +x /usr/local/bin/${BINARY_NAME}
	cp ./server-shutdown.service /etc/systemd/system/server-shutdown.service