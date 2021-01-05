module github.com/BlueMedoraPublic/terraform-provider-bindplane

go 1.14

require (
	github.com/BlueMedoraPublic/bpcli v1.3.0
	github.com/google/uuid v1.1.1
	github.com/hashicorp/terraform v0.12.25
	github.com/pkg/errors v0.9.1
	golang.org/x/sys v0.0.0-20190922100055-0a153f010e69 // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
