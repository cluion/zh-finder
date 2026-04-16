package matcher

import (
	"github.com/cluion/zh-finder/internal/classifier"
)

type Match struct {
	Text  string
	Type  classifier.HanType
	Start int
	End   int
}

type Matcher struct {
	classifier *classifier.Classifier
}

func New(c *classifier.Classifier) *Matcher {
	return &Matcher{classifier: c}
}

func (m *Matcher) Find(content string) []Match {
	var results []Match
	runes := []rune(content)
	i := 0

	for i < len(runes) {
		if classifier.IsHan(runes[i]) {
			start := i
			var text []rune
			typeCounts := make(map[classifier.HanType]int)

			for i < len(runes) && classifier.IsHan(runes[i]) {
				text = append(text, runes[i])
				t := m.classifier.Classify(runes[i])
				typeCounts[t]++
				i++
			}

			mainType := m.determineMainType(typeCounts)
			results = append(results, Match{
				Text:  string(text),
				Type:  mainType,
				Start: start,
				End:i,
			})
		} else {
			i++
		}
	}

	return results
}

func (m *Matcher) determineMainType(counts map[classifier.HanType]int) classifier.HanType {
	if counts[classifier.Traditional] > 0 && counts[classifier.Simplified] > 0 {
		return classifier.Common
	}
	if counts[classifier.Traditional] > counts[classifier.Simplified] {
		return classifier.Traditional
	}
	if counts[classifier.Simplified] > 0 {
		return classifier.Simplified
	}
	return classifier.Common
}
