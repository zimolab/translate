package translate

import (
	"embed"
	"path/filepath"
	"testing"
)

const TestDataDir = "test_data"

//go:embed test_data
var BuiltinTranslations embed.FS

func TestNewTranslator1(t *testing.T) {
	defaultTranslationFile := filepath.Join(TestDataDir, "active.en_US.toml")
	tr, err := NewTranslator("active", nil, defaultTranslationFile)
	if err != nil {
		t.Error(err)
	}
	if tr.Tr("ID_TEST", "error") != "hello, world!" {
		t.Fail()
	}
}

func TestNewTranslator2(t *testing.T) {
	// 测试无法创建对象的情形
	// parseLanguage()失败：不符合文件名模式
	defaultTranslationFile := filepath.Join(TestDataDir, "translate.zh-CN.toml")
	_, err := NewTranslator("active", nil, defaultTranslationFile)
	if err == nil {
		t.Fail()
	} else {
		t.Log("PASS", err)
	}
	// parseLanguage()失败：language-tag不符合BCP 47规范
	defaultTranslationFile = filepath.Join(TestDataDir, "active.mylang.toml")
	_, err = NewTranslator("active", nil, defaultTranslationFile)
	if err == nil {
		t.Fail()
	} else {
		t.Log("PASS", err)
	}
	// 文件不存在
	defaultTranslationFile = filepath.Join(TestDataDir, "active.en-US.toml")
	_, err = NewTranslator("active", nil, defaultTranslationFile)
	if err == nil {
		t.Fail()
	} else {
		t.Log("PASS", err)
	}
}

func TestNewTranslator3(t *testing.T) {
	// 测试从内嵌文件系统加载翻译文件
	defaultTranslationFile := "test_data/active.zh_CN.toml"
	_, err := NewTranslator("active", BuiltinTranslations, defaultTranslationFile)
	if err != nil {
		t.Error(err)
	}
}

func TestTranslator_LoadTranslationFile(t *testing.T) {
	defaultTranslationFile := filepath.Join(TestDataDir, "active.en_US.toml")
	tr, err := NewTranslator("active", nil, defaultTranslationFile)
	if err != nil {
		t.Error(err)
	}

	err = tr.LoadTranslationFile(nil, filepath.Join(TestDataDir, "active.zh_CN.toml"))
	if err != nil {
		t.Error(err)
	}

	err = tr.LoadTranslationFile(BuiltinTranslations, "test_data/active.zh-CN.toml")
	if err != nil {
		t.Error(err)
	}

	trs := tr.AllTranslations(true)
	if len(trs) != 3 {
		t.Fail()
	} else {
		t.Log("translations:", trs)
		t.Log("tags:", tr.bundle.LanguageTags())
	}

}

func TestTranslator_SetTranslation(t *testing.T) {
	// 创建实例，并将active.en_US.toml加载为默认翻译文件
	defaultTranslationFile := filepath.Join(TestDataDir, "active.en_US.toml")
	tr, err := NewTranslator("active", nil, defaultTranslationFile)
	if err != nil {
		t.Error(err)
	}
	// 加载active.zh_CN.toml
	err = tr.LoadTranslationFile(nil, filepath.Join(TestDataDir, "active.zh_CN.toml"))
	if err != nil {
		t.Error(err)
	}
	// 获取ID_TEST对应文本
	if tr.Tr("ID_TEST", "error") != "hello, world!" {
		t.Log("logic error")
		t.Fail()
	}
	// 切换翻译语言
	err = tr.SetTranslation("简体中文")
	if err != nil {
		t.Error(err)
	}
	// 再次获取ID_TEST对应文本
	if tr.Tr("ID_TEST", "error") != "你好，世界！" {
		t.Log("logic error")
		t.Fail()
	}
	// 切换未被加载的语言
	err = tr.SetTranslation("中文")
	if err == nil {
		t.Log("logic error")
		t.Fail()
	}
	// 测试翻译文件覆盖效果
	// zh-CN、zh_CN会被视为同一种语言
	// 因此后加载的active.zh-CN.toml会覆盖前面加载的active.zh_CN.toml
	err = tr.LoadTranslationFile(nil, filepath.Join(TestDataDir, "active.zh-CN.toml"))
	if err != nil {
		t.Error(err)
	}
	if tr.Tr("ID_TEST", "error") != "世界，你好！" {
		t.Log("logic error")
		t.Fail()
	}
	msg1 := tr.Tr("ID_TEST", "error1")
	// 且此时“简体中文”、“中文”两个displayName都对应了同一个zh—CN标签
	err = tr.SetTranslation("中文")
	if err != nil {
		t.Error(err)
	}
	msg2 := tr.Tr("ID_TEST", "error2")
	if msg1 != msg2 {
		t.Log("logic error")
		t.Fail()
	}
}
