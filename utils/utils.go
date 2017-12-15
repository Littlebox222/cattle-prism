package utils

import (
	"regexp"
	"strconv"
)

func IdStringToIdNumber(a string) int64 {

	if a == "" {
		return 0
	}

	reg := regexp.MustCompile(`\d+`)
	itemIdNums := reg.FindAllStringSubmatch(a, -1)

	if len(itemIdNums) != 2 {
		return 0
	}

	if itemIdNum, err := strconv.ParseInt(itemIdNums[1][0], 10, 64); err != nil {
		return 0
	} else {
		return itemIdNum
	}
}
