package evaluator

import (
	"fmt"

	"github.com/spaolacci/murmur3"
	"github.com/tuannguyensn2001/aurora-go/auroratype"
)

func CalculateHash(value interface{}, key string) uint32 {
	valueStr := fmt.Sprintf("%v", value)
	combinedStr := key + ":" + valueStr
	return murmur3.Sum32([]byte(combinedStr))
}

func IsInPercentageRange(hash uint32, percentage int) bool {
	if percentage <= 0 {
		return false
	}
	if percentage >= 100 {
		return true
	}

	const numBuckets = 10000
	hashBucket := int(hash % numBuckets)
	threshold := percentage * (numBuckets / 100)

	return hashBucket < threshold
}

func SelectVariantByHash(experimentID string, hashAttribute string, attr map[string]any, variants []auroratype.Variant) *auroratype.Variant {
	if len(variants) == 0 {
		return nil
	}

	if len(variants) == 1 {
		return &variants[0]
	}

	totalRollout := 0
	for _, v := range variants {
		totalRollout += v.Rollout
	}

	hashValue := attr[hashAttribute]
	if hashValue == nil {
		return nil
	}

	hash := CalculateHash(hashValue, experimentID+":"+hashAttribute)

	const numBuckets = 10000
	hashBucket := int(hash % numBuckets)

	normalizedBucket := (hashBucket * totalRollout) / numBuckets

	cumulative := 0
	for _, v := range variants {
		cumulative += v.Rollout
		if normalizedBucket < cumulative {
			return &v
		}
	}

	return &variants[len(variants)-1]
}
