package utils

import "regexp"

// FindStringNamedSubmatches returns a map of the named capturing groups of the provided
// regexp in the given string.
// Based on the example at http://blog.kamilkisiel.net/blog/2012/07/05/using-the-go-regexp-package/.
func FindStringNamedSubmatches(r *regexp.Regexp, s string) map[string]string {

	submatches := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return submatches
	}

	for i, name := range r.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}

		submatches[name] = match[i]

	}
	return submatches
}
