BINARIES = servus-extapi servus-player servus-tragedy servus-boardgames servus-magissa
GOOPTS = -buildmode=pie -trimpath -mod=readonly -modcacherw -ldflags=-s -ldflags=-w

build: format $(BINARIES)

format:
	gofmt -s -w .

$(BINARIES):
	CGO_ENABLED=0 go build $(GOOPTS) -o dist/$@ ./cmd/$@
