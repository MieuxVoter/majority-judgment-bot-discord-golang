package locales

import (
	"embed"
	"github.com/BurntSushi/toml"
	"github.com/disgoorg/disgo/discord"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"log"
	"main/src/container"
	"main/src/locales/guesser"
)

//go:embed locale.*.toml
var LocaleFS embed.FS

func GetLocalization() *Localization {
	return container.Get("localization").(*Localization)
}

func GetLocalizer(languages ...string) *Localizer {
	return GetLocalization().GetLocalizer(languages...)
}

func GetServerLocalizer() *Localizer {
	return GetLocalizer(
		guesser.DetectLanguagesFromEnv(language.AmericanEnglish)...,
	)
}

type Localizer struct {
	Localizer *i18n.Localizer
}

// T translates the message identified by its key
func (l *Localizer) T(key string) string {
	s, _ := l.Localizer.LocalizeMessage(&i18n.Message{ID: key})
	return s
}

// Tf translates and formats the message identified by its key
func (l *Localizer) Tf(key string, data map[string]interface{}) string {
	s, _ := l.Localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: key},
		TemplateData:   data,
	})
	return s
}

// Tp translates and pluralizes the message identified by its key
func (l *Localizer) Tp(
	key string,
	amount interface{},
) string {
	s, _ := l.Localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: key},
		PluralCount:    amount,
	})
	return s
}

// Tfp translates, formats and pluralizes the message identified by its key
func (l *Localizer) Tfp(
	key string,
	amount interface{},
	data map[string]interface{},
) string {
	s, _ := l.Localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: key},
		TemplateData:   data,
		PluralCount:    amount,
	})
	return s
}

type Localization struct {
	logger *logrus.Logger
	bundle *i18n.Bundle
}

func (l *Localization) Init() {
	l.bundle = i18n.NewBundle(language.AmericanEnglish)
	l.bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	dirEntries, err := LocaleFS.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		_, err = l.bundle.LoadMessageFileFS(LocaleFS, dirEntry.Name())
		if err != nil {
			panic(err)
		}
	}
}

func (l *Localization) GetLocalizer(languages ...string) *Localizer {
	return &Localizer{Localizer: i18n.NewLocalizer(l.bundle, languages...)}
}

// GetLanguages returns the available languages, starting with the default one.
func (l *Localization) GetLanguages() []string {
	languages := make([]string, 0)
	for _, tag := range l.bundle.LanguageTags() {
		languages = append(languages, tag.String())
	}
	return languages
}

func (l *Localization) GetTranslations(
	key string,
) map[discord.Locale]string {
	all := make(map[discord.Locale]string, 0)

	for _, lang := range l.GetLanguages() {
		localizer := l.GetLocalizer(lang)
		locale := discord.Locale(lang)
		localized := localizer.T(key)
		if localized != "" {
			all[locale] = localized
		}
	}

	return all
}

func init() {
	err := container.GetBuilder().Add(di.Def{
		Name: "localization",
		Build: func(ctn di.Container) (interface{}, error) {
			service := &Localization{
				logger: ctn.Get("logger").(*logrus.Logger),
			}
			service.Init()
			return service, nil
		},
	})
	if err != nil {
		log.Fatalln("localization failed to build:", err)
	}
}
