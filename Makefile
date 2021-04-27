GOMOD=mdv
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
HASH=$(shell git log -n1 --pretty=format:%h)
REVS=$(shell git log --oneline|wc -l)
build: debug
upx:
	upx -9 $(GOMOD)*
debug: setver compdbg
release: setver comprel upx
linux: export GOOS=linux
linux: export GOARCH=amd64
linux: release
windows: export GOOS=windows
windows: export GOARCH=amd64
windows: release
mac: export GOOS=darwin
mac: export GOARCH=amd64
mac: release
setver:
	sed 's/{_BRANCH}/$(BRANCH)/' verinfo.tpl > version.sed.1
	sed 's/{_G_HASH}/$(HASH)/' version.sed.1 > version.sed.2
	sed 's/{_G_REVS}/$(REVS)/' version.sed.2 > version.go
	rm -fr version.sed*
comprel:
	go build -ldflags="-s -w" .
compdbg:
	go build -race -gcflags=all=-d=checkptr=0 .
clean:
	rm -fr $(GOMOD)* version.go
