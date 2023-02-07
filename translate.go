package translate

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
)

// 翻译文件的文件名模式：<prefix>.<lang-tag>.toml
var filenamePattern = `%s\.([\w-]+)\.toml`

type displayname struct {
	Display string
}

type Translator struct {
	// 翻译文件前缀
	prefix string
	// 翻译文件名匹配模式
	regex *regexp.Regexp
	// i18n bundle
	bundle *i18n.Bundle
	// language tag映射表：显示名 -> tag
	languages map[string]string
	// 当前Localizer
	localizer *i18n.Localizer
}

// NewTranslator
//
//	@Description: 创建Translator对象
//	@param prefix 翻译文件名前缀
//	@param fsys fs.FS，为nil则代表从磁盘读取文件
//	@param defaultFile 默认翻译文件
//	@return *Translator
//	@return error
func NewTranslator(prefix string, fsys fs.FS, defaultFile string) (*Translator, error) {

	t := &Translator{}

	t.prefix = prefix
	t.languages = map[string]string{}
	t.regex = regexp.MustCompile(fmt.Sprintf(filenamePattern, prefix))

	// 从默认翻译文件的文件名中解析langTag、tagName
	langTag, tagName, err := t.parseLanguageTag(defaultFile)
	if err != nil {
		return nil, err
	}

	// 创建bundle
	t.bundle = i18n.NewBundle(*langTag)
	// 注册反序列化方法，使用toml作为翻译文件格式
	t.bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 加载默认翻译文件
	err = t.LoadTranslationFile(fsys, defaultFile)
	if err != nil {
		return nil, err
	}
	// 创建初始localizer
	t.localizer = i18n.NewLocalizer(t.bundle, tagName)
	return t, nil
}

// 从文件名中解析语言标签，成功解析需要满足以下两点：
// 1.文件名的形式符合<prefix>.<language_tag>.toml的格式，换言之，可以被fileRegExp正则表达式匹配
// 2.language_tag本身符合BCP 47规则，即能够通过language.Parse()方法的解析
func (t *Translator) parseLanguageTag(filename string) (langTag *language.Tag, tagName string, err error) {
	// 利用正则表达式，匹配出文件名中的language tagName
	// 如<prefix>.zh_CN.toml => tagName == "zh_CN"
	// 如<prefix>.en-US.toml => tagName == "en-US"
	result := t.regex.FindStringSubmatch(filename)
	if len(result) < 2 {
		return nil, "", illegalFilename(filename)
	}
	tagName = result[1]
	if tagName == "" {
		return nil, "", illegalFilename(filename)
	}
	// 解析该tag，确保其tag符合BCP 47等规则
	tmp, err := language.Parse(tagName)
	if err != nil {
		return nil, "", err
	}
	if &tmp == nil {
		return nil, "", illegalFilename(filename)
	}
	// 返回解析出的language.Tag对象、tag名称
	langTag = &tmp
	return langTag, tagName, nil
}

// 从翻译文件（toml文件）中解析该翻译文件对应语言的显示名称
// 该名称由Display字段定义，如：
// # <prefix>.zh-CN.toml
// Display = "简体中文"
func (t *Translator) parseDisplayName(fsys fs.FS, path string) (string, error) {
	m := displayname{}
	var err error
	if fsys == nil {
		_, err = toml.DecodeFile(path, &m)
	} else {
		_, err = toml.DecodeFS(fsys, path, &m)
	}
	if err != nil {
		return "", err
	}
	return m.Display, nil
}

// LoadTranslationFile
//
//	@Description: 加载翻译文件，并解析其语言标签和显示名
//	@receiver t
//	@param fsys 文件系统，若为nil，则从磁盘读取
//	@param path 文件路径
//	@return error
func (t *Translator) LoadTranslationFile(fsys fs.FS, path string) error {
	filename := filepath.Base(path)
	// 解析tag
	_, tag, err := t.parseLanguageTag(filename)
	if err != nil {
		return err
	}
	// 解析显示名
	displayName, err := t.parseDisplayName(fsys, path)
	if err != nil {
		return err
	}
	if displayName == "" {
		displayName = tag
	}
	// 加载文件
	if fsys == nil {
		_, err = t.bundle.LoadMessageFile(path)
	} else {
		_, err = t.bundle.LoadMessageFileFS(fsys, path)
	}
	if err != nil {
		return err
	}
	// 文件加载成功，将displayName与tag关联起来
	t.languages[displayName] = tag
	return nil
}

// SetTranslation
//
//	@Description: 设置当前翻译
//	@receiver t
//	@param displayName
//	@return error
func (t *Translator) SetTranslation(displayName string) error {
	// 获取显示名所对应tag
	tagName, ok := t.languages[displayName]
	if !ok {
		return notLoaded(displayName)
	}
	t.localizer = i18n.NewLocalizer(t.bundle, tagName)
	return nil
}

func (t *Translator) AllTranslations(shouldSort bool) []string {
	if len(t.languages) == 0 {
		return nil
	}
	keys := make([]string, 0, len(t.languages))
	for _, key := range t.languages {
		keys = append(keys, key)
	}
	if shouldSort {
		// 排序后再返回，确保显示顺序的一致性
		sort.Strings(keys)
	}
	return keys
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

func (t *Translator) Tr(messageId string, fallback string) string {
	return t.Translate(messageId, fallback)
}

func (t *Translator) Localizer() *i18n.Localizer {
	return t.localizer
}
