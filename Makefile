GOMOD=mdv
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
HASH=$(shell git log -n1 --pretty=format:%h)
REVS=$(shell git log --oneline|wc -l)
debug: setver compdbg
release: setver comprel
linux: export GOOS=linux
linux: export GOARCH=amd64
linux: release
windows: export GOOS=windows
windows: export GOARCH=amd64
windows: release
mac: export GOOS=darwin
mac: export GOARCH=amd64
mac: release
dist:
	make
	tar czf mdv.$(REVS).$(HASH)-linux-amd64.tar.gz mdv
	rm mdv
	make mac
	tar czf mdv.$(REVS).$(HASH)-darwin-amd64.tar.gz mdv
	rm mdv
	make windows
	zip mdv.$(REVS).$(HASH)-windows-amd64.zip mdv.exe
	rm mdv.exe
setver:
	sed 's/{_BRANCH}/$(BRANCH)/' verinfo.tpl > version.sed.1
	sed 's/{_G_HASH}/$(HASH)/' version.sed.1 > version.sed.2
	sed 's/{_G_REVS}/$(REVS)/' version.sed.2 > version.go
	rm -fr version.sed*
comprel:
	CGO_ENABLED=0 go build -ldflags="-s -w" .
compdbg:
	go build -race -gcflags=all=-d=checkptr=0 .
clean:
	rm -fr $(GOMOD)* version.go
distclean: clean