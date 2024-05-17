package tools

import (
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	files, _ := os.ReadDir("./tools/lang")
	for _, file := range files {
		if file.Type().IsRegular() && filepath.Ext(file.Name()) == ".yaml" {
			bundle.MustLoadMessageFile("./tools/lang/" + file.Name())
		}
	}
}

func I18n(idioma language.Tag, mensaje string) string {
	localizer := i18n.NewLocalizer(bundle, idioma.String())
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    mensaje,
			Other: mensaje,
		},
	})
}
