# staging environment retrieves dependencies and compiles
#
FROM golang:1.15

WORKDIR /terraform-provider-bindplane

ARG version

RUN \
    apt-get update >> /dev/null && \
    apt-get install -y zip unzip wget

ADD . /terraform-provider-bindplane

RUN go test ./...

# compile with gox
RUN go get github.com/mitchellh/gox

# Disable CGO to avoid pulling in C dependencies, and compile for
# MACOS, Linux, and Windows
RUN \
    env CGO_ENABLED=0 \
    $GOPATH/bin/gox \
        -arch=amd64 \
        -os='!netbsd !openbsd !freebsd'  \
        -output "artifacts/terraform-provider-bindplane_{{.OS}}_{{.Arch}}" \
        -tags '-a' \
        -ldflags '-w -extldflags "-static"' \
        ./...

# install terraform
RUN \
    wget https://releases.hashicorp.com/terraform/0.12.24/terraform_0.12.24_linux_amd64.zip && \
    unzip terraform_0.12.24_linux_amd64.zip && \
    mv terraform /usr/bin

# install terraform-provider-bindplane provider and overwrite
# any possible installs already in place
WORKDIR /terraform-provider-bindplane/artifacts
RUN yes | cp -rf terraform-provider-bindplane_linux_amd64 /terraform-provider-bindplane/example/basic/terraform-provider-bindplane

# smoke test: make sure init and validate works
WORKDIR /terraform-provider-bindplane/example/basic
RUN /usr/bin/terraform init
RUN /usr/bin/terraform validate

# rename each binary and then zip them
WORKDIR /terraform-provider-bindplane/artifacts
RUN mv terraform-provider-bindplane_linux_amd64 terraform-provider-bindplane_v${version} && zip terraform-provider-bindplane_linux_amd64_v${version}.zip terraform-provider-bindplane_v${version}
RUN mv terraform-provider-bindplane_darwin_amd64 terraform-provider-bindplane_v${version} && zip terraform-provider-bindplane_darwin_amd64_v${version}.zip terraform-provider-bindplane_v${version}
RUN mv terraform-provider-bindplane_windows_amd64.exe terraform-provider-bindplane_v${version}.exe && zip terraform-provider-bindplane_windows_amd64_v${version}.zip terraform-provider-bindplane_v${version}.exe

# build the sha256sum file
RUN ls | grep 'terraform-provider-bindplane_' | xargs -n1 sha256sum >> SHA256SUMS

# keep only the zip files
RUN ls | grep -Ev 'zip|SUM' | xargs -n1 rm -f
