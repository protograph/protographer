package protographer

import (
	"fmt"
	"regexp"
	"strings"
)

// normalize trim the string, remove whitespace, and convert it to lowercase.
func normalize(l, p string) string {
	return strings.ToLower(strings.TrimSpace(strings.TrimPrefix(l, p)))
}

func expandLine(l string, expr [][2]string) string {

	var re *regexp.Regexp
	// if it is a block quote, change it to math mode.
	if strings.HasPrefix(l, "  ") {
		l = "$" + strings.TrimSpace(l) + "$"
	}

	// find stuff in math mode
	re = regexp.MustCompile("[$][^$]*[$]")
	l = re.ReplaceAllStringFunc(l, func(s string) string {

		// from $foo$ to foo
		s = s[1 : len(s)-1]

		const boundaryStart = "(^|[ -.:-@^_{}-~])"
		const boundaryEnd = "($|[ -.:-@^_{}-~])"

		// find stuff that is not between `...`
		// the trick is, add ` before and after s, then extract anything with `...`.
		re2 := regexp.MustCompile("[`][^`]*[`]")
		s = re2.ReplaceAllStringFunc("`"+s+"`", func(t string) string {
			// from `bar` to bar
			t = t[1 : len(t)-1]
			// find english words, with boundary with special character (except :,[,\,])
			re3 := regexp.MustCompile(boundaryStart + "([A-Za-z]{3,})")
			t = re3.ReplaceAllString(t, "$1`$2`")
			return "`" + t + "`"
		})
		s = s[1 : len(s)-1]

		// apply expressions
		for _, e := range expr {
			re4 := regexp.MustCompile(boundaryStart + e[0] + boundaryEnd)
			s = re4.ReplaceAllString(s, "${1}"+e[1]+"${2}")
		}

		// other transformation
		s = strings.Replace(s, "'", `\textnormal{\textquotesingle}`, -1)
		s = strings.Replace(s, "<-", `\gets`, -1)
		s = strings.Replace(s, "->", `\to`, -1)
		s = strings.Replace(s, "||", `\Vert `, -1)

		return "$" + s + "$"
	})

	// expand `foo` to textop (with proper spacing)
	re = regexp.MustCompile(" `([^`]*)` ")
	l = re.ReplaceAllString(l, "{\\textop{ $1 }}")
	re = regexp.MustCompile(" `([^`]*)`")
	l = re.ReplaceAllString(l, "{\\textop{ $1}}")
	re = regexp.MustCompile("`([^`]*)` ")
	l = re.ReplaceAllString(l, "{\\textop{$1 }}")
	re = regexp.MustCompile("`([^`]*)`")
	l = re.ReplaceAllString(l, "{\\textop{$1}}")

	// expand <- to \gets
	l = strings.Replace(l, "<-", "\\gets", -1)
	return l
}

func l2r(a bool, f int) string {
	if a && f == 0 {
		return "east"
	}
	if !a && f == 1 {
		return "east"
	}
	return "west"
}

func expand(a string) string {
	a = strings.TrimSpace(a)
	// expand EOL to \n and shortstack
	if strings.Contains(a, "\n") {
		a = fmt.Sprintf("\\shortstack[l]{%s}", strings.Replace(a, "\n", "\\\\", -1))
	}
	return a
}

func padtop(j int) string {
	s := ""
	for i := 0; i < j; i++ {
		s += "\\postlevel\n"
	}
	return s
}
