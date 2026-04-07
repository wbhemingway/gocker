package cli

import (
	"fmt"
	"strconv"
	"strings"
)

func parseRate(input string) (int64, error) {
	if input == "" {
		return 0, nil
	}

	parts := strings.Split(input, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid format: too many decimals")
	}

	dollars, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || dollars < 0 || strings.HasPrefix(parts[0], "-") {
		return 0, fmt.Errorf("invalid dollar amount")
	}

	cents := int64(0)
	if len(parts) == 2 {
		centStr := parts[1]
		if len(centStr) > 2 {
			return 0, fmt.Errorf("cannot have more than two decimal places")
		}
		if len(centStr) == 1 {
			centStr += "0"
		}

		cents, err = strconv.ParseInt(centStr, 10, 64)
		if err != nil || cents < 0 {
			return 0, fmt.Errorf("invalid cent amount")
		}
	}

	return (dollars * 100) + cents, nil
}

func parseTags(input string) []string {
	if input == "" {
		return nil
	}
	var tagList []string
	for tag := range strings.SplitSeq(input, ",") {
		cleanTag := strings.TrimSpace(tag)
		if cleanTag != "" {
			tagList = append(tagList, cleanTag)
		}
	}
	return tagList
}
