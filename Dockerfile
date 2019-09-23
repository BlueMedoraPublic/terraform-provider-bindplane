# staging environment retrieves dependencies and compiles
#
FROM golang:1.12

WORKDIR /src/bindplane-terraform-provider

ENV GOPATH=/src

ARG version

RUN \
    apt-get update >> /dev/null && \
    apt-get install -y golint zip

ADD . /src/bindplane-terraform-provider

# compile with gox
RUN go get github.com/mitchellh/gox

# Disable CGO to avoid pulling in C dependencies, and compile for
# MACOS, Linux, and Windows
RUN \
    env CGO_ENABLED=0 \
    $GOPATH/bin/gox \
        -arch=amd64 \
        -os='!netbsd !openbsd !freebsd'  \
        -tags '-a' \
        -ldflags '-w -extldflags "-static"' \
        ./...

# rename each binary and then zip them
RUN mv terraform-provider-bindplane_linux_amd64 terraform-provider-bindplane && zip terraform-provider-bindplane_linux_amd64.zip terraform-provider-bindplane
RUN mv terraform-provider-bindplane_darwin_amd64 terraform-provider-bindplane && zip terraform-provider-bindplane_darwin_amd64.zip terraform-provider-bindplane
RUN mv terraform-provider-bindplane_windows_amd64.exe terraform-provider-bindplane.exe && zip terraform-provider-bindplane_windows_amd64.zip terraform-provider-bindplane.exe

# build the sha256sum file
RUN ls | grep 'terraform-provider-bindplane_' | xargs -n1 sha256sum >> SHA256SUMS
