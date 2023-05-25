package boot

import "strings"

// 返回值childParma，patternParam，ok
func dynamicMatch(childPattern, pattern string) (string, string, bool) {
	if childPattern[len(childPattern)-1] != ':' {
		return "", "", false
	}
	cps := strings.Split(childPattern, ":")
	switch len(cps) {
	case 3:
		ps := strings.Split(pattern, ":")
		if ps[0] == cps[0] {
			return cps[1], ps[1], true
		}
	case 2:
		if strings.Contains(pattern, ":") {
			return "", "", false
		}
		return cps[0], pattern, true
	}
	return "", "", false
}
