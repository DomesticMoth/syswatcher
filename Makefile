dependencies:
	go get
	
build: dependencies main.go
	go build -o build/syswatcher main.go

install:
	cp build/syswatcher /usr/bin/syswatcher
	mkdir -p /etc/syswatcher/
	cp default_config.toml /etc/syswatcher/config.toml
	cp syswatcher.service /lib/systemd/system/syswatcher.service

clear:
	rm -rf build/
