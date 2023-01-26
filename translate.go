package translate

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/fs"
	"path/filepath"
	"regexp"
)

type _meta struct {
	Display string
}

type Translator struct {
	translateFilePrefix string
	translateFileRegExp *regexp.Regexp
	bundle              *i18n.Bundle
	tagMap              map[string]string
	currentTag          string
	localizer           *i18n.Localizer
}

func NewTranslator(translateFilePrefix string, defaultTranslateFile string) (*Translator, error) {

	t, err := createTranslator(translateFilePrefix, defaultTranslateFile)
	if err != nil {
		return nil, err
	}

	displayName, err := t.getDisplayName(defaultTranslateFile)
	if err != nil {
		return nil, err
	}
	if displayName == "" {
		displayName = t.currentTag
	}
	t.tagMap[displayName] = t.currentTag

	_, err = t.bundle.LoadMessageFile(defaultTranslateFile)
	if err != nil {
		return nil, err
	}

	t.localizer = i18n.NewLocalizer(t.bundle, t.currentTag)

	return t, nil
}

func NewTranslatorFS(translateFilePrefix string, fs fs.FS, defaultTranslateFile string) (*Translator, error) {

	t, err := createTranslator(translateFilePrefix, defaultTranslateFile)
	if err != nil {
		return nil, err
	}

	displayName, err := t.getDisplayNameFS(fs, defaultTranslateFile)
	if err != nil {
		return nil, err
	}
	if displayName == "" {
		displayName = t.currentTag
	}
	t.tagMap[displayName] = t.currentTag

	_, err = t.bundle.LoadMessageFile(defaultTranslateFile)
	if err != nil {
		return nil, err
	}

	t.localizer = i18n.NewLocalizer(t.bundle, t.currentTag)

	return t, nil
}

func createTranslator(translateFilePrefix string, defaultTranslateFile string) (*Translator, error) {
	t := &Translator{}

	t.translateFilePrefix = translateFilePrefix
	t.tagMap = map[string]string{}

	t.translateFileRegExp = regexp.MustCompile(fmt.Sprintf(`%s\.([\w-]+)\.toml`, translateFilePrefix))

	langTag := t.getLangTag(defaultTranslateFile)
	if langTag == "" {
		return nil, illegalFilename(langTag)
	}

	lang, err := language.Parse(langTag)
	if err != nil {
		return nil, illegalLanguageTag(langTag)
	}

	t.currentTag = langTag
	t.bundle = i18n.NewBundle(lang)
	t.bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	return t, nil
}

func (t *Translator) getLangTag(filename string) string {
	tag := t.translateFileRegExp.FindStringSubmatch(filename)
	if tag == nil || len(tag) <= 1 {
		return ""
	}
	return tag[1]
}

func (t *Translator) getDisplayName(path string) (string, error) {
	meta := _meta{}
	_, err := toml.DecodeFile(path, &meta)
	if err != nil {
		return "", err
	}
	return meta.Display, nil
}

func (t *Translator) getDisplayNameFS(fs fs.FS, path string) (string, error) {
	meta := _meta{}
	_, err := toml.DecodeFS(fs, path, &meta)
	if err != nil {
		return "", err
	}
	return meta.Display, nil
}

func (t *Translator) IsLanguageTagLegal(tag string) bool {
	_, err := language.Parse(tag)
	return err == nil
}

func (t *Translator) LoadTranslateFile(path string) error {
	filename := filepath.Base(path)
	tag := t.getLangTag(filename)
	if tag == "" {
		return illegalFilename(filename)
	}

	if !t.IsLanguageTagLegal(tag) {
		return illegalLanguageTag(tag)
	}

	displayName, err := t.getDisplayName(path)
	if err != nil {
		return err
	}
	if displayName == "" {
		displayName = tag
	}

	_, err = t.bundle.LoadMessageFile(path)
	if err != nil {
		return err
	}

	t.tagMap[displayName] = tag

	return nil
}

func (t *Translator) LoadTranslateFileFS(fs fs.FS, path string) error {
	filename := filepath.Base(path)
	tag := t.getLangTag(filename)
	if tag == "" {
		return illegalFilename(filename)
	}

	if !t.IsLanguageTagLegal(tag) {
		return illegalLanguageTag(tag)
	}

	displayName, err := t.getDisplayNameFS(fs, path)
	if err != nil {
		return err
	}
	if displayName == "" {
		displayName = tag
	}

	_, err = t.bundle.LoadMessageFileFS(fs, path)
	if err != nil {
		return err
	}

	t.tagMap[displayName] = tag

	return nil
}

func (t *Translator) LanguageTagOf(displayName string) string {
	tag, ok := t.tagMap[displayName]
	if ok {
		return tag
	}
	return ""
}

func (t *Translator) SetCurrentTranslation(displayName string) error {
	tag := t.LanguageTagOf(displayName)
	if tag == "" {
		return translationNotLoaded(displayName)
	}
	t.currentTag = tag
	t.localizer = i18n.NewLocalizer(t.bundle, tag)
	return nil
}

func (t *Translator) Translate(messageId string, fallback string) string {
	localized, err := t.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageId,
		DefaultMessage: &i18n.Message{
			ID:    messageId,
			Other: fallback,
		},
	})
	if err != nil {
		return fallback
	} else {
		return localized
	}
}

func (t *Translator) GetCurrentLocalizer() *i18n.Localizer {
	return t.localizer
}
