package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var onlyWords = regexp.MustCompile(`[^\s\x{2000}-\x{206F}\x{2E00}-\x{2E7F}\\'!"#$%&()*+,./:;<=>?@\[\]^_{|}~]+`)

type textObj struct {
	count int
	val   string
}

func Top10(text string) []string {
	toCount := onlyWords.FindAllString(text, -1)
	if len(text) == 0 || len(toCount) == 0 {
		return nil
	}
	mapOfStr := make(map[string]*textObj, len(toCount))
	for _, s := range toCount {
		key := strings.ToLower(s)
		if key == "-" {
			continue
		}
		if obj, ok := mapOfStr[key]; ok {
			obj.count++
		} else {
			mapOfStr[key] = &textObj{count: 1, val: key}
		}
	}
	sliToSort := make([]textObj, 0, len(toCount))
	for i := range mapOfStr {
		sliToSort = append(sliToSort, *mapOfStr[i])
	}
	sort.Slice(sliToSort, func(i, j int) bool {
		if sliToSort[j].count == sliToSort[i].count {
			return sliToSort[j].val > sliToSort[i].val
		}
		return sliToSort[j].count < sliToSort[i].count
	})
	out := make([]string, 10)
	for i := 0; i < 10; i++ {
		out[i] = sliToSort[i].val
	}
	return out
}
