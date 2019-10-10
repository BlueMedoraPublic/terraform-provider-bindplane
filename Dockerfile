# staging environment retrieves dependencies and compiles
#
FROM golang:1.12

WORKDIR /src/terraform-provider-bindplane

ARG version

RUN \
    apt-get update >> /dev/null && \
    apt-get install -y golint zip unzip wget

ADD . /src/terraform-provider-bindplane

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

# install terraform
RUN \
    wget https://releases.hashicorp.com/terraform/0.12.10/terraform_0.12.10_linux_amd64.zip && \
    unzip terraform_0.12.10_linux_amd64.zip && \
    mv terraform /usr/bin

# install terraform-provider-bindplane provider
RUN cp terraform-provider-bindplane_linux_amd64 /usr/bin/terraform-provider-bindplane_v${version}

# smoke test: make sure init and validate works
WORKDIR /src/terraform-provider-bindplane/example/basic
RUN /usr/bin/terraform init
RUN /usr/bin/terraform validate

# rename each binary and then zip them
WORKDIR /src/terraform-provider-bindplane
RUN mv terraform-provider-bindplane_linux_amd64 terraform-provider-bindplane_v${version} && zip terraform-provider-bindplane_linux_amd64_v${version}.zip terraform-provider-bindplane_v${version}
RUN mv terraform-provider-bindplane_darwin_amd64 terraform-provider-bindplane_v${version} && zip terraform-provider-bindplane_darwin_amd64_v${version}.zip terraform-provider-bindplane_v${version}
RUN mv terraform-provider-bindplane_windows_amd64.exe terraform-provider-bindplane_v${version}.exe && zip terraform-provider-bindplane_windows_amd64_v${version}.zip terraform-provider-bindplane_v${version}.exe

# build the sha256sum file
RUN ls | grep 'terraform-provider-bindplane_' | xargs -n1 sha256sum >> SHA256SUMS
