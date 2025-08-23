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
	servus-series \
	servus-auth \
	database-exporter \
	file-exporter \
	servus-deco \
	servus-network \
	servus-people \
	servus-bgg \
	servus-hg8010h \
	servus-speedportplus2

GOOPTS = -buildmode=pie -trimpath -mod=readonly -modcacherw -ldflags=-s -ldflags=-w
OUTPUT_DIR = dist

build: format $(BINARIES)

format:
	gofmt -s -w .

generate-values-files: $(addprefix $(OUTPUT_DIR)/, $(addsuffix .values.yaml, $(BINARIES)))

clean:
	rm $(OUTPUT_DIR)/*.yaml

$(BINARIES):
	CGO_ENABLED=0 go build $(GOOPTS) -o $(OUTPUT_DIR)/$@ ./cmd/$@

package-%:
	ko resolve -f $(OUTPUT_DIR)/$*.values.yaml --platform=linux/arm64 --sbom=none --preserve-import-paths > charts/$*/image.yaml

update-version-%: package-%
	@VERSION=$$(grep 'version": "' ./cmd/$*/main.go | awk -F'"' '{print $$4}'); \
	echo "Updating version to $$VERSION for $*"; \
	yq eval -i ".appVersion = \"$$VERSION\"" ./charts/$*/Chart.yaml

deploy-%: update-version-% $(OUTPUT_DIR)/%.values.yaml
	helm upgrade --install $* ./charts/$* -f charts/$*/image.yaml

$(OUTPUT_DIR)/%.values.yaml:
	@yq -n '.image.pullPolicy = "IfNotPresent"' > $@
	@yq -i '.image.tag = "ko://github.com/DictumMortuum/servus-extapi/cmd/$*"' $@
	@yq -i '.imagePullSecrets = []' $@
