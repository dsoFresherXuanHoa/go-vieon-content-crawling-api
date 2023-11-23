package utils

import "golang.org/x/exp/slices"

var ()

type sliceUtil struct{}

func NewSliceUtil() *sliceUtil {
	return &sliceUtil{}
}

func (sliceUtil) RemoveDuplicatedStringSlice(slice []string, defaultRemove []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, item := range slice {
		if _, ok := keys[item]; !ok && !slices.Contains(defaultRemove, item) {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func (sliceUtil) Diff(firstSlice []string, secondSlice []string) []string {
	result := make([]string, 0)
	for _, i := range secondSlice {
		found := false
		for _, j := range firstSlice {
			if i == j {
				found = true
				break
			}
		}
		if !found {
			result = append(result, i)
		}
	}
	return result
}
