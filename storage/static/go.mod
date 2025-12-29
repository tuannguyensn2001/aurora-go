module github.com/tuannguyensn2001/aurora-go/storage/static

go 1.25.5

require (
	github.com/tuannguyensn2001/aurora-go/core v0.0.0
	gopkg.in/yaml.v2 v2.4.0
)

require github.com/spaolacci/murmur3 v1.1.0 // indirect

replace github.com/tuannguyensn2001/aurora-go => ../../core
