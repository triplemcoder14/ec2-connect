package helpers

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func StrOrDefault(s *string, defaultVal string) string {
	if s == nil {
		return defaultVal
	}
	return *s
}

func GetTagName(inst *types.Instance) string {
	var nameValue string
	for _, t := range inst.Tags {
		if *t.Key == "Name" {
			nameValue = *t.Value
		}
	}
	return nameValue
}

func filter(slice []string, predicate func(string) bool) []string {
	var result []string
	for _, elem := range slice {
		if predicate(elem) {
			result = append(result, elem)
		}
	}
	return result
}

func contains(slice []string, elem string) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}
