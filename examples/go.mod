module examples

go 1.25.5

require (
	github.com/tuannguyensn2001/aurora-go v0.0.0-00010101000000-000000000000
	github.com/tuannguyensn2001/aurora-go/storage/static v0.0.0-00010101000000-000000000000
)

require (
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/tuannguyensn2001/aurora-go/auroratype v0.0.0-00010101000000-000000000000 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/tuannguyensn2001/aurora-go => ../

replace github.com/tuannguyensn2001/aurora-go/storage/static => ../storage/static

replace github.com/tuannguyensn2001/aurora-go/auroratype => ../auroratype
