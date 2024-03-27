all:
	cp ./configs/defaultConf.yaml /etc/lmuConf.yaml
	go mod init linux-monitoring-utility
	go mod tidy
	go build -o ./build/lmu ./cmd/main.go 


clean:
	rm -r build
	rm go.mod