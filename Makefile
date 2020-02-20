all:
	go build .
	mkdir -p ${HOME}/go/src/github.com/yonasadiel/helios
	cp -r ./* ${HOME}/go/src/github.com/yonasadiel/helios
