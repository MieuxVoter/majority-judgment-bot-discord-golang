package locales

import (
	"embed"
	"github.com/BurntSushi/toml"
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
		guesser.DetectLanguagesFromEnv(language.English)...,
	)
}

type Localizer struct {
	Localizer *i18n.Localizer
}

func (l *Localizer) T(key string) string {
	s, _ := l.Localizer.LocalizeMessage(&i18n.Message{ID: key})
	return s
}

type Localization struct {
	logger *logrus.Logger
	bundle *i18n.Bundle
}

func (l *Localization) Init() {
	l.bundle = i18n.NewBundle(language.English)
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
