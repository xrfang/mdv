GOMOD=mdv
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
HASH=$(shell git log -n1 --pretty=format:%h)
REVS=$(shell git log --oneline|wc -l)
build: debug
upx:
	upx -9 $(GOMOD)*
debug: setver compdbg
release: setver comprel upx
windows: export GOOS=windows
windows: export GOARCH=amd64
windows: release
mac: export GOOS=darwin
mac: export GOARCH=amd64
mac: release
setver:
	cp verinfo.tpl version.go
	sed -i 's/{_BRANCH}/$(BRANCH)/' version.go
	sed -i 's/{_G_HASH}/$(HASH)/' version.go
	sed -i 's/{_G_REVS}/$(REVS)/' version.go
comprel:
	go build -ldflags="-s -w" .
compdbg:
	go build -race -gcflags=all=-d=checkptr=0 .
clean:
	rm -fr $(GOMOD)* version.go
