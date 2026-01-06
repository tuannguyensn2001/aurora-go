module github.com/tuannguyensn2001/aurora-go/fetcher/s3

go 1.25.5

require (
	github.com/aws/aws-sdk-go-v2 v1.36.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.76.1
	github.com/tuannguyensn2001/aurora-go/auroratype v0.0.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.8 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.32 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.32 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.32 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.6.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.13 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
)

replace github.com/tuannguyensn2001/aurora-go/auroratype => ../../auroratype
