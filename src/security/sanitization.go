package security

import (
	"regexp"
	"strings"
)

const NoLimit = -1

var delimiter = regexp.MustCompile("[|]")

// regexp2 is nice because it supports lookbehind and lookahead, but for now it has no Split(),
// and implementing a Split should be done upstream here : https://github.com/dlclark/regexp2
//var delimiter = regexp2.MustCompile("(?<![|])[|](?![|])", regexp2.None)

func ExtractProposalsNames(rawName string) []string {
	rawNames := delimiter.Split(rawName, NoLimit)
	// De-double pipes (only works well with negative lookbehind regex, which we don't have for now)
	//for k := range rawNames {
	//	rawNames[k] = strings.Replace(rawNames[k], "||", "|", NoLimit)
	//}
	return rawNames
}

var markdownFontStyle = regexp.MustCompile("[*_`]")

func RemoveMarkdown(raw string) string {
	return strings.TrimSpace(markdownFontStyle.ReplaceAllString(raw, " "))
}

func EscapeCsvValue(value string) string {
	return strings.ReplaceAll(value, "\"", "")
}

// TruncateString truncates a string value to a maximum length
func TruncateString(value string, maxLength int) string {
	actualLength := len(value)
	if actualLength > maxLength {
		actualLength = maxLength
	}

	return value[:actualLength]
}

// TruncateEllipsis returns a truncated string with a trailing ellipsis (…) if relevant
func TruncateEllipsis(value string, maxLength int) string {
	if maxLength == 0 {
		return ""
	}
	actualLength := len(value)
	if actualLength <= maxLength {
		return value
	}
	if actualLength < 2 {
		return value
	}
	if maxLength == 1 {
		return value[:1]
	}
	actualLength = maxLength - 1
	return value[:actualLength] + "…"
}
