package translate

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/emirpasic/gods/maps/treebidimap"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

// 翻译文件的文件名模式：<prefix>.<lang-tag>.toml
const localeFilenamePattern = `%s\.([\w-]+)\.toml`

type _LocaleFileMeta struct {
	DisplayName string `toml:"displayname"`
}

type Locales struct {
	regex     *regexp.Regexp
	bundle    *i18n.Bundle
	languages *treebidimap.Map
	localizer *i18n.Localizer
}

func NewLocales(localeFilenamePrefix string, defaultLanguage string) (*Locales, error) {
	l := &Locales{
		regex:     regexp.MustCompile(f(localeFilenamePattern, localeFilenamePrefix)),
		languages: treebidimap.NewWithStringComparators(),
		localizer: nil,
	}
	tag, err := language.Parse(defaultLanguage)
	if err != nil {
		return nil, err
	}
	l.bundle = i18n.NewBundle(tag)
	l.bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	return l, err

}

func (l *Locales) LoadLocaleFile(fsys fs.FS, path string) (displayName, tagName string, err error) {
	// 解析文件
	tag, meta, err := l.parseFile(fsys, path)
	if err != nil {
		return "", "", err
	}
	tagName = tag.String()
	displayName = meta.DisplayName
	// 解析成功后，加载文件
	err = l.loadFile(fsys, path)
	if err != nil {
		return displayName, tagName, err
	}
	// 文件加载成功后，把DisplayName和tagName保存到languages表中
	// 查询tagName是否被加载过
	oldName, ok := l.languages.GetKey(tagName)
	if ok {
		// 若已加载过，则用新的DisplayName进行替换
		l.languages.Remove(oldName)
		l.languages.Put(displayName, tagName)
	} else {
		// 若未加载过，则建立DisplayName到tagName的对应
		l.languages.Put(displayName, tagName)
	}
	return displayName, tagName, nil
}

func (l *Locales) LoadLocalesDir(dir string) (success []string, fail []string, err error) {
	allTomlFiles := make([]string, 0, 20)
	// 先找出目录下所有toml文件
	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".toml" {
			allTomlFiles = append(allTomlFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	// 逐一读取
	for _, path := range allTomlFiles {
		_, _, err1 := l.LoadLocaleFile(nil, path)
		if err1 != nil {
			fail = append(fail, path)
		} else {
			success = append(success, path)
		}
	}
	return success, fail, nil
}

func (l *Locales) SetLocaleByName(displayName string) (tagName string, err error) {
	tmp, ok := l.languages.Get(displayName)
	if !ok {
		return "", notLoaded(displayName)
	}
	tagName = tmp.(string)
	l.localizer = i18n.NewLocalizer(l.bundle, tagName)
	return tagName, nil
}

func (l *Locales) SetLocaleByTag(tagName string) (displayName string, err error) {
	tmp, ok := l.languages.GetKey(tagName)
	if !ok {
		return "", notLoaded(tagName)
	}
	displayName = tmp.(string)
	l.localizer = i18n.NewLocalizer(l.bundle, tagName)
	return displayName, nil
}

func (l *Locales) TagNameOf(displayName string) (tagName string, err error) {
	tmp, ok := l.languages.Get(displayName)
	if !ok {
		return "", notLoaded(displayName)
	}
	tagName = tmp.(string)
	return tagName, err
}

func (l *Locales) DisplayNameOf(tagName string) (displayName string, err error) {
	tmp, ok := l.languages.GetKey(tagName)
	if !ok {
		return "", notLoaded(tagName)
	}
	displayName = tmp.(string)
	return displayName, err
}

func (l *Locales) GetLocales() (tagNames []string) {
	tagNames = make([]string, 0, l.languages.Size())
	for _, tagName := range l.languages.Values() {
		tagNames = append(tagNames, tagName.(string))
	}
	return tagNames
}

func (l *Locales) GetLocaleNames() (displayNames []string) {
	displayNames = make([]string, 0, l.languages.Size())
	for _, displayName := range l.languages.Keys() {
		displayNames = append(displayNames, displayName.(string))
	}
	return displayNames
}

func (l *Locales) Tr(id string) (string, error) {
	return l.Localize(&i18n.LocalizeConfig{
		MessageID: id,
	})
}

func (l *Locales) Localize(config *i18n.LocalizeConfig) (string, error) {
	if l.localizer == nil {
		return "", localeNotSet()
	}
	return l.localizer.Localize(config)
}

func (l *Locales) loadFile(fsys fs.FS, path string) error {
	var err error
	if fsys == nil {
		_, err = l.bundle.LoadMessageFile(path)
	} else {
		_, err = l.bundle.LoadMessageFileFS(fsys, path)
	}
	return err
}

func (l *Locales) parseFile(fsys fs.FS, path string) (language.Tag, _LocaleFileMeta, error) {
	// 解析tag
	filename := filepath.Base(path)
	tagPart, err := l.parseTagPart(filename)
	if err != nil {
		return language.Tag{}, _LocaleFileMeta{}, err
	}
	tag, err := language.Parse(tagPart)
	if err != nil {
		return tag, _LocaleFileMeta{}, err
	}
	// 解析meta
	meta, err := parseMeta(fsys, path)
	if err != nil {
		return tag, meta, err
	}
	return tag, meta, nil
}

func (l *Locales) parseTagPart(filename string) (string, error) {
	// 利用正则表达式，匹配出文件名中的语言标签部分（即tagPart）
	// 如<prefix>.zh_CN.toml => tagPart == "zh_CN"
	// 如<prefix>.en-US.toml => tagPart == "en-US"
	result := l.regex.FindStringSubmatch(filename)
	if len(result) < 2 {
		return "", illegalFilename(filename)
	}
	tagPart := strings.TrimSpace(result[1])
	if tagPart == "" {
		return tagPart, illegalFilename(filename)
	}
	return tagPart, nil
}

func parseMeta(fsys fs.FS, path string) (_LocaleFileMeta, error) {
	m := _LocaleFileMeta{
		DisplayName: "",
	}
	var err error
	if fsys == nil {
		_, err = toml.DecodeFile(path, &m)
	} else {
		_, err = toml.DecodeFS(fsys, path, &m)
	}
	if err != nil {
		return m, err
	}
	m.DisplayName = strings.TrimSpace(m.DisplayName)
	// 确保displayname字段不为空，为空则返回一个error
	if m.DisplayName == "" {
		return m, missingDisplayName(path)
	}
	return m, nil
}

func f(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}
