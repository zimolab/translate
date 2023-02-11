package translate

import (
	"embed"
	"testing"
)

//go:embed test_data
var embedFs embed.FS

func TestNewLocales(t *testing.T) {
	locales, err := NewLocales("active", "zh-CN")
	if err != nil || locales == nil {
		t.Error(err)
	}
}

func TestNewLocales_2(t *testing.T) {
	locales, err := NewLocales("active", "zh_CN")
	if err != nil || locales == nil {
		t.Error(err)
	}
}

func TestNewLocales_3(t *testing.T) {
	locales, err := NewLocales("active", "zh")
	if err != nil || locales == nil {
		t.Error(err)
	}
}

func TestNewLocales_4(t *testing.T) {
	locales, err := NewLocales("active", "foo")
	if err == nil || locales != nil {
		t.Log("Logic Error, should produce an error here")
	}
	t.Log("illegal language-tag produces an error: ", err)
}

func TestLocales_LoadLocaleFile_1(t *testing.T) {
	t.Log("Load a locale file")
	locales, _ := NewLocales("active", "art")
	displayName, tagName, err := locales.LoadLocaleFile(nil, "test_data/active.en_US.toml")
	if err != nil {
		t.Error(err)
	}
	t.Log("Loaded: ", displayName, "<=>", tagName)
}

func TestLocales_LoadLocaleFile_2(t *testing.T) {
	t.Log("Load locale files")
	locales, _ := NewLocales("active", "art")
	_, _, err := locales.LoadLocaleFile(nil, "test_data/active.en_US.toml")
	if err != nil {
		t.Error(err)
	}
	_, _, err = locales.LoadLocaleFile(nil, "test_data/active.zh-CN.toml")
	if err != nil {
		t.Error(err)
	}
	t.Log("Loaded:", locales.languages.String())
}

func TestLocales_LoadLocaleFile_3(t *testing.T) {
	t.Log("Load locale files with same display names")
	locales, _ := NewLocales("active", "art")
	_, _, err := locales.LoadLocaleFile(nil, "test_data/active.en_US.toml")
	if err != nil {
		t.Error(err)
	}
	t.Log("Before: ", locales.languages)
	_, _, err = locales.LoadLocaleFile(nil, "test_data/active.en_AU.toml")
	if err != nil {
		t.Error(err)
	}
	tags := locales.GetLocales()
	displaynames := locales.GetLocaleNames()
	if len(tags) != 1 || tags[0] != "en-AU" || len(displaynames) != 1 || displaynames[0] != "English" {
		t.Error("Logic Error")
	}
	t.Log("After: ", locales.languages)
}

func TestLocales_LoadLocaleFile_4(t *testing.T) {
	t.Log("Load locale files with same tag name")
	locales, _ := NewLocales("active", "art")
	_, _, err := locales.LoadLocaleFile(nil, "test_data/active.en_US.toml")
	if err != nil {
		t.Error(err)
	}
	t.Log("Before: ", locales.languages)

	_, _, err = locales.LoadLocaleFile(nil, "test_data/active.en-US.toml")
	if err != nil {
		t.Error(err)
	}
	tags := locales.GetLocales()
	displaynames := locales.GetLocaleNames()
	if len(tags) != 1 || tags[0] != "en-US" || len(displaynames) != 1 || displaynames[0] != "English(US)" {
		t.Error("Logic Error")
		t.Fail()
	}
	t.Log("After: ", locales.languages)
}

func TestLocales_LoadLocaleFile_5(t *testing.T) {
	t.Log("Load locale from embed fs")
	locales, _ := NewLocales("active", "art")
	_, _, _ = locales.LoadLocaleFile(embedFs, "test_data/active.en-US.toml")
	_, _, _ = locales.LoadLocaleFile(embedFs, "test_data/active.zh-CN.toml")
	tags := locales.GetLocales()
	t.Log("Loaded: ", tags)
	if tags[0] != "en-US" || tags[1] != "zh-CN" {
		t.Error("Logic Error")
	}
}

func TestLocales_LoadLocaleFile_6(t *testing.T) {
	t.Log("Load locale with errors")
	locales, _ := NewLocales("active", "art")
	_, _, err := locales.LoadLocaleFile(embedFs, "test_data/active.badtag.toml")
	if err == nil {
		t.Error("Logic Error")
	} else {
		t.Log("a not well-formed tag should produce an error: ", err)
	}

	_, _, err = locales.LoadLocaleFile(embedFs, "test_data/active.foo.toml")
	if err == nil {
		t.Error("Logic Error")
	} else {
		t.Log("an unknown tag should produce an error too: ", err)
	}

	_, _, err = locales.LoadLocaleFile(embedFs, "test_data/active.en.toml")
	if err == nil {
		t.Error("Logic Error")
	} else {
		t.Log("illegal displayname should also produce an error: ", err)
	}

	_, _, err = locales.LoadLocaleFile(embedFs, "test_data/translate.zh-CN.toml")
	if err == nil {
		t.Error("Logic Error")
	} else {
		t.Log("bad locale filename should also produce an error: ", err)
	}

	_, _, err = locales.LoadLocaleFile(embedFs, "test_data/active.af.toml")
	if err == nil {
		t.Error("Logic Error")
	} else {
		t.Log("bad toml format should also produce an error: ", err)
	}
	_, _, err = locales.LoadLocaleFile(embedFs, "test_data/active.af-AZ.toml")
	if err == nil {
		t.Error("Logic Error")
	} else {
		t.Log("bad field type should also produce an error: ", err)
	}

	t.Log("tags: ", locales.GetLocales())
}

func TestLocales_TagNameOf(t *testing.T) {
	locales, _ := NewLocales("active", "art")
	_, _, _ = locales.LoadLocaleFile(embedFs, "test_data/active.en-US.toml")
	_, _, _ = locales.LoadLocaleFile(embedFs, "test_data/active.zh-CN.toml")
	t.Log("Loaded: ", locales.languages)
	tagName1, _ := locales.TagNameOf("中文")
	tagName2, _ := locales.TagNameOf("English(US)")
	if tagName1 != "zh-CN" || tagName2 != "en-US" {
		t.Error("Logic Error")
	}
	_, err := locales.TagNameOf("简体中文")
	if err == nil {
		t.Error("Logic Error")
	} else {
		t.Log("get the tag name of an unloaded locale with its displayname should produce an error: ", err)
	}
}

func TestLocales_DisplayNameOf(t *testing.T) {
	locales, _ := NewLocales("active", "art")
	_, _, _ = locales.LoadLocaleFile(embedFs, "test_data/active.en-US.toml")
	_, _, _ = locales.LoadLocaleFile(embedFs, "test_data/active.zh-CN.toml")
	t.Log("Loaded: ", locales.languages)
	displayName1, _ := locales.DisplayNameOf("zh-CN")
	displayName2, _ := locales.DisplayNameOf("en-US")
	if displayName1 != "中文" || displayName2 != "English(US)" {
		t.Error("Logic Error")
	}
	_, err := locales.DisplayNameOf("af-AZ")
	if err == nil {
		t.Error("Logic Error")
	} else {
		t.Log("get the display name of an unloaded locale with its tag name should produce an error: ", err)
	}
}

func TestLocales_SetLocaleByName(t *testing.T) {
	locales, _ := NewLocales("active", "art")
	_, _, _ = locales.LoadLocaleFile(embedFs, "test_data/active.en-US.toml")
	_, _, _ = locales.LoadLocaleFile(embedFs, "test_data/active.zh-CN.toml")
	t.Log("Loaded: ", locales.languages)
	_, err := locales.Tr("ID_TEST")
	if err == nil {
		t.Error("Logic Error")
	} else {
		t.Log("invoke locales.Tr() before locales.SetLocaleBy**() should produce an error: ", err)
		t.Log("\n")
	}
	tagName, _ := locales.SetLocaleByName("中文")
	t.Log("current tag: ", tagName)
	msg, _ := locales.Tr("ID_TEST")
	if msg != "世界，你好！" {
		t.Error("Logic Error")
	}
	t.Log("translated: ", msg)

	tagName, _ = locales.SetLocaleByName("English(US)")
	t.Log("current tag: ", tagName)
	msg, _ = locales.Tr("ID_TEST")
	if msg != "hello, world!" {
		t.Error("Logic Error")
	}
	t.Log("translated: ", msg)

	tagName, err = locales.SetLocaleByName("English")
	if err == nil {
		t.Error("Logic Error")
	} else {
		t.Log("locales.SetLocaleByName() with an unloaded locale should produce an error: ", err)
	}
}

func TestLocales_SetLocaleByTag(t *testing.T) {
	locales, _ := NewLocales("active", "art")
	_, _, _ = locales.LoadLocaleFile(embedFs, "test_data/active.en-US.toml")
	_, _, _ = locales.LoadLocaleFile(embedFs, "test_data/active.zh-CN.toml")
	t.Log("Loaded: ", locales.languages)
	_, err := locales.Tr("ID_TEST")
	if err == nil {
		t.Error("Logic Error")
	} else {
		t.Log("invoke locales.Tr() before locales.SetLocaleBy**() should produce an error: ", err)
	}
	displayName, _ := locales.SetLocaleByTag("zh-CN")
	t.Log("current tag: ", displayName)
	msg, _ := locales.Tr("ID_TEST")
	if msg != "世界，你好！" {
		t.Error("Logic Error")
	}
	t.Log("translated: ", msg)

	displayName, _ = locales.SetLocaleByTag("en-US")
	t.Log("current tag: ", displayName)
	msg, _ = locales.Tr("ID_TEST")
	if msg != "hello, world!" {
		t.Error("Logic Error")
	}
	t.Log("translated: ", msg)

	displayName, err = locales.SetLocaleByTag("en")
	if err == nil {
		t.Error("Logic Error")
	} else {
		t.Log("locales.SetLocaleByTag() with an unloaded locale should produce an error: ", err)
	}
}

func TestLocales_parseTagPart(t *testing.T) {
	locales, _ := NewLocales("active", "art")
	_, err := locales.parseTagPart("active.  .toml")
	if err == nil {
		t.Error("Logic Error")
	} else {
		t.Log(err)
	}
}

func TestLocales_LoadLocalesDir(t *testing.T) {
	locales, _ := NewLocales("active", "art")
	success, fail, _ := locales.LoadLocalesDir("test_data")
	t.Log("success: ", success)
	t.Log("fail: ", fail)
}

func TestLocales_LoadLocalesDir2(t *testing.T) {
	locales, _ := NewLocales("active", "art")
	success, fail, err := locales.LoadLocalesDir("test_data1")
	t.Log("success: ", success)
	t.Log("fail: ", fail)
	t.Log(err)
}
