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
