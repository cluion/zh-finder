package classifier

type HanType int

const (
	AllHanType = iota
	Traditional
	Simplified
	Common
)

func (t HanType) String() string {
	switch t {
	case Traditional:
		return "traditional"
	case Simplified:
		return "simplified"
	default:
		return "common"
	}
}

type Classifier struct {
	charMap map[rune]HanType
}

func New(traditionalData, simplifiedData string) *Classifier {
	c := &Classifier{
		charMap: make(map[rune]HanType),
	}
	c.loadCharacterLists(traditionalData, simplifiedData)
	return c
}

func (c *Classifier) loadCharacterLists(tradData, simpData string) {
	tradChars := c.loadCharsFromString(tradData)
	simpChars := c.loadCharsFromString(simpData)

	// Traditional only
	for ch := range tradChars {
		if !simpChars[ch] {
			c.charMap[ch] = Traditional
		}
	}

	// Simplified only
	for ch := range simpChars {
		if !tradChars[ch] {
			c.charMap[ch] = Simplified
		}
	}

	// Common
	for ch := range tradChars {
		if simpChars[ch] {
			c.charMap[ch] = Common
		}
	}
}

func (c *Classifier) loadCharsFromString(data string) map[rune]bool {
	chars := make(map[rune]bool)
	for _, r := range data {
		if r != '\n' && r != '\r' && r != ' ' {
			chars[r] = true
		}
	}
	return chars
}

func (c *Classifier) Classify(ch rune) HanType {
	if !IsHan(ch) {
		return Common
	}
	if t, ok := c.charMap[ch]; ok {
		return t
	}
	return Common
}

func IsHan(ch rune) bool {
	// CJK Unified Ideographs
	return (ch >= 0x4E00 && ch <= 0x9FFF) ||
		// CJK Unified Ideographs Extension A
		(ch >= 0x3400 && ch <= 0x4DBF) ||
		// CJK Unified Ideographs Extension B-F
		(ch >= 0x20000 && ch <= 0x2CEAF) ||
		// CJK Compatibility Ideographs
		(ch >= 0xF900 && ch <= 0xFAFF)
}
