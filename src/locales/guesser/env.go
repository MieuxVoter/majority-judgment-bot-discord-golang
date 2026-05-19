package guesser

import (
	"golang.org/x/text/language"
	"os"
)

// envVariablesHoldingLocale is sorted by decreasing priority (breaks on first found)
// These environment variables are expected to hold a parsable locale (fr_FR, es, en-US, …)
// ADR: https://www.gnu.org/software/gettext/manual/html_node/Locale-Environment-Variables.html
var envVariablesHoldingLocale = []string{
	"LANGUAGE",
	"LC_ALL",
	"LANG",
}

func DetectLanguagesFromEnv(defaultLanguage language.Tag) []string {
	var detectedLanguages []string
	for _, envKey := range envVariablesHoldingLocale {
		lang := os.Getenv(envKey)
		if lang != "" {
			detectedLang := language.Make(lang)
			appendLang(&detectedLanguages, detectedLang)
		}
	}
	appendLang(&detectedLanguages, defaultLanguage)

	return detectedLanguages
}

func appendLang(languages *[]string, lang language.Tag) {
	langString := lang.String()
	*languages = append(*languages, langString)

	langBase, confidentInBase := lang.Base()
	if confidentInBase != language.No {
		*languages = append(*languages, langBase.String())
		*languages = append(*languages, langBase.ISO3())
	}
}
