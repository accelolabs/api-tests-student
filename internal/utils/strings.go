package utils

import "unicode"

func MixCase(text string) string {
	mixed := []rune(text)
	upper := true

	for i, r := range mixed {
		if !unicode.IsLetter(r) {
			continue
		}

		upper = !upper

		if upper {
			mixed[i] = unicode.ToUpper(r)
			continue
		}

		mixed[i] = unicode.ToLower(r)
	}

	return string(mixed)
}
