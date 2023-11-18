BINARIES = \
	servus-extapi \
	servus-player \
	servus-tragedy \
	servus-boardgames \
	servus-modem-a \
	servus-modem-restart-a \
	servus-modem-b \
	servus-modem-c \
	servus-modem-d \
	servus-bgstats \
	servus-prices \
	servus-scrape \
	servus-series

GOOPTS = -buildmode=pie -trimpath -mod=readonly -modcacherw -ldflags=-s -ldflags=-w

build: format $(BINARIES)

format:
	gofmt -s -w .

$(BINARIES):
	CGO_ENABLED=0 go build $(GOOPTS) -o dist/$@ ./cmd/$@
