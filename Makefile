ifeq (, $(shell which docker))
    $(error "No docker in $(PATH)")
endif

VERSION := $(shell cat main.go | grep "const version" | cut -c 17- | tr -d '"')

$(shell mkdir -p artifacts)

build: clean fmt
	$(info building terraform-provider-bindplane ${VERSION})

	@docker build \
		--no-cache \
	    --build-arg version=${VERSION} \
	    -t btp:${VERSION} .

	@docker create -ti --name artifacts btp:${VERSION} bash && \
		docker cp artifacts:/terraform-provider-bindplane/artifacts/. artifacts/

	# cleanup
	@docker rm -fv artifacts &> /dev/null

test:
	go test ./...

lint:
	golint ./...

fmt:
	go fmt ./...

quick:
	$(shell env CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"')

clean:
	$(shell rm -rf artifacts/*)
