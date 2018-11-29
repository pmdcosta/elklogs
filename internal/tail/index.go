package tail

import (
	"fmt"
	"regexp"
	"time"

	"github.com/pkg/errors"
)

func FilterIndex(indices []string, indexPattern string, start *time.Time, end *time.Time) ([]string, error) {
	// check if the query is date filtered
	if start == nil && end == nil {
		index := findLastIndex(indices, indexPattern)
		return []string{index}, nil
	}

	if start == nil {
		start = &time.Time{}
	}
	if end == nil {
		t := time.Now()
		end = &t
	}

	result := make([]string, 0, len(indices))
	for _, idx := range indices {
		matched, _ := regexp.MatchString(indexPattern, idx)
		if matched {
			idxDate, err := extractIndexDate(idx)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed parsing log date: %s", idx))
			}
			if (idxDate.After(*start) || idxDate.Equal(*start)) && (idxDate.Before(*end) || idxDate.Equal(*end)) {
				result = append(result, idx)
			}
		}
	}
	return result, nil
}

// findLastIndex sorts the indices that match the pattern and returns the most recent one
func findLastIndex(indices []string, indexPattern string) string {
	var lastIdx string
	for _, idx := range indices {
		matched, _ := regexp.MatchString(indexPattern, idx)
		if matched {
			if &lastIdx == nil {
				lastIdx = idx
			} else if idx > lastIdx {
				lastIdx = idx
			}
		}
	}
	return lastIdx
}

// extractIndexDate extracts and parses the index date from its name
func extractIndexDate(dateStr string) (*time.Time, error) {
	dateRegexp := regexp.MustCompile(fmt.Sprintf(`(\d{4}.\d{2}.\d{2})`))
	match := dateRegexp.FindAllStringSubmatch(dateStr, -1)
	if len(match) == 0 {
		return nil, fmt.Errorf("failed to extract date: %s", dateStr)
	}
	result := match[0]
	parsed, err := time.Parse("2006.01.02", result[0])
	if err != nil {
		return nil, errors.Wrap(err, "failed parsing date")
	}
	return &parsed, nil
}
