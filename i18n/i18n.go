package i18n

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	go18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"

	"github.com/Lumicrate/gompose/http"
)

const (
	CtxLanguages = "x-languages"
	CtxLocalizer = "x-localizer"
)

type LanguageExtractorOptions map[string]any
type LanguageExtractor func(http.Context, LanguageExtractorOptions) []string

type Translator struct {
	locale                   string
	defaultLanguage          string
	languageExtractors       []LanguageExtractor
	languageExtractorOptions LanguageExtractorOptions
	bundle                   *go18n.Bundle
}

func NewI18n(directory, defaultLanguage string) (*Translator, error) {
	files, err := listFiles(directory)
	if err != nil {
		return nil, err
	}

	t := &Translator{
		defaultLanguage: defaultLanguage,
		languageExtractors: []LanguageExtractor{
			CookieLanguageExtractor,
			HeaderLanguageExtractor,
		},
		languageExtractorOptions: LanguageExtractorOptions{
			"CookieName":    "lang",
			"SessionName":   "lang",
			"URLPrefixName": "lang",
		},
		bundle: go18n.NewBundle(language.Make(defaultLanguage)),
		locale: defaultLanguage,
	}

	t.bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	for _, file := range files {
		data, err := readFile(directory, file)
		if err != nil {
			log.Printf("i18n: read file %v return error: %v\r\n", file, err.Error())
			continue
		}

		if _, err := t.bundle.ParseMessageFileBytes(data, file); err != nil {
			fmt.Printf("i18n: parse message file %v return error: %v\r\n", file, err.Error())
			continue
		}
	}

	return t, nil
}

func listFiles(directory string) ([]string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), "yaml") || strings.HasSuffix(entry.Name(), "yml") {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func readFile(directory, filename string) ([]byte, error) {
	path := filepath.Join(directory, filename)
	return os.ReadFile(path)
}

func CookieLanguageExtractor(c http.Context, o LanguageExtractorOptions) []string {
	langs := make([]string, 0)
	if cookieName := o["CookieName"].(string); cookieName != "" {
		if cookie, err := c.Request().Cookie(cookieName); err == nil {
			if cookie.Value != "" {
				langs = append(langs, cookie.Value)
			}
		}
	} else {
		log.Println("i18n middleware: \"CookieName\" is not defined in LanguageExtractorOptions")
	}
	return langs
}

func HeaderLanguageExtractor(c http.Context, o LanguageExtractorOptions) []string {
	langs := make([]string, 0)
	acceptLang := c.Request().Header.Get("Accept-Language")
	if acceptLang != "" {
		langs = append(langs, parseAcceptLanguage(acceptLang)...)
	}
	return langs
}

func URLPrefixLanguageExtractor(c http.Context, o LanguageExtractorOptions) []string {
	langs := make([]string, 0)
	if urlPrefixName := o["URLPrefixName"].(string); urlPrefixName != "" {
		paramLang := c.Param(urlPrefixName)
		if paramLang != "" && strings.HasPrefix(c.Request().URL.Path, fmt.Sprintf("/%s", paramLang)) {
			langs = append(langs, paramLang)
		}
	} else {
		log.Println("i18n middleware: \"URLPrefixName\" is not defined in LanguageExtractorOptions")
	}
	return langs
}

func parseAcceptLanguage(acptLang string) []string {
	var lqs []string

	langQStrs := strings.Split(acptLang, ",")
	for _, langQStr := range langQStrs {
		trimedLangQStr := strings.Trim(langQStr, " ")

		langQ := strings.Split(trimedLangQStr, ";")
		lq := langQ[0]
		lqs = append(lqs, lq)
	}
	return lqs
}

func localize(localizer *go18n.Localizer, messageID string, data map[string]any) string {
	msg, err := localizer.Localize(&go18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})

	if err != nil {
		return messageID
	}
	return msg
}

func (t *Translator) extractLanguages(c http.Context) []string {
	langs := make([]string, 0)
	for _, extractor := range t.languageExtractors {
		langs = append(langs, extractor(c, t.languageExtractorOptions)...)
	}

	langs = append(langs, t.defaultLanguage)
	return langs
}

func (t *Translator) GetMiddleware() http.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(ctx http.Context) {
			if langs := ctx.Get(CtxLanguages); langs == nil {
				ctx.Set(CtxLanguages, t.extractLanguages(ctx))
			}

			if localizer := ctx.Get(CtxLocalizer); localizer == nil {
				langs := ctx.Get(CtxLanguages).([]string)

				localizer = go18n.NewLocalizer(t.bundle, langs...)

				ctx.Set(CtxLocalizer, localizer)
			}

			next(ctx)
		}
	}
}

func (t *Translator) AddMessage(locale, id, message string) (err error) {
	msg := &go18n.Message{
		ID:    id,
		Other: message,
	}

	if err = t.bundle.AddMessages(language.Make(locale), msg); err != nil {
		return err
	}

	return nil
}

func (t *Translator) T(messageID string, args ...any) string {
	var (
		data      map[string]any
		lang      string
		localizer *go18n.Localizer
	)

	for _, arg := range args {
		switch v := arg.(type) {
		case map[string]any:
			data = v
		case string:
			lang = v
		}
	}

	if lang == "" {
		lang = t.locale
	}

	localizer = go18n.NewLocalizer(t.bundle, lang)

	return localize(localizer, messageID, data)
}

func (t *Translator) SetLocale(locale string) *Translator {
	t.locale = locale

	return t
}
