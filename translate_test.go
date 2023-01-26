package translate

import (
	"embed"
	"testing"
)

//go:embed active.*.toml
var LocaleFS embed.FS

func TestNewTranslator1(t *testing.T) {
	trans, err := NewTranslator("active", "a.en.toml")
	if trans != nil || err == nil {
		t.Fail()
	} else {
		t.Log("err:", err)
	}
}

func TestNewTranslator2(t *testing.T) {
	trans, err := NewTranslator("active", "active.en.toml")
	if trans != nil || err == nil {
		t.Fail()
	} else {
		t.Log("err:", err)
	}
}

func TestNewTranslator3(t *testing.T) {
	trans, err := NewTranslator("active", "active.foo.toml")
	if trans != nil || err == nil {
		t.Fail()
	} else {
		t.Log("err:", err)
	}
}

func TestNewTranslator4(t *testing.T) {
	trans, err := NewTranslator("active", "test_data/active.zh.toml")
	if trans != nil || err == nil {
		t.Fail()
	} else {
		t.Log("err:", err)
	}
}

func TestNewTranslator5(t *testing.T) {
	trans, err := NewTranslator("active", "test_data/active.zh_CN.toml")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if trans.currentTag != "zh_CN" {
		t.Fail()
	}

	t.Log(trans.tagMap)
}

func TestNewTranslator6(t *testing.T) {
	trans, err := NewTranslator("active", "test_data/active.zh-CN.toml")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if trans.currentTag != "zh-CN" {
		t.Fail()
	}

	t.Log(trans.tagMap)
}

func TestTranslator_LoadTranslateFile(t *testing.T) {
	trans, err := NewTranslator("active", "./test_data/active.zh_CN.toml")
	if err != nil {
		t.Fatal(err)
	}
	err = trans.LoadTranslateFile("active.toml")
	if err == nil {
		t.Fail()
	}
	err = trans.LoadTranslateFile("active.foo.toml")
	if err == nil {
		t.Fail()
	}
	err = trans.LoadTranslateFile("active.en_US.toml")
	if err != nil {
		t.Fail()
	}
	err = trans.LoadTranslateFile("./test_data/active.en_US.toml")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestTranslator_SetCurrentTranslation(t *testing.T) {
	trans, err := NewTranslator("active", "./test_data/active.zh_CN.toml")
	if err != nil {
		t.Fatal(err)
	}
	err = trans.LoadTranslateFile("./test_data/active.en_US.toml")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(trans.tagMap)
	t.Log(trans.localizer)
	t.Log(trans.currentTag)

	err = trans.SetCurrentTranslation("简体中文")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(trans.localizer)
	t.Log(trans.currentTag)

	err = trans.SetCurrentTranslation("en_US")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(trans.localizer)
	t.Log(trans.currentTag)

	err = trans.SetCurrentTranslation("en-US")
	if err == nil {
		t.Error(err)
		t.Fail()
	} else {
		t.Log(err.Error())
	}
	t.Log(trans.localizer)
	t.Log(trans.currentTag)

}

func TestFSFunc(t *testing.T) {
	transFS, err := NewTranslatorFS("active", LocaleFS, "active.zh_CN.toml")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(transFS.tagMap)
	t.Log(transFS.currentTag)
	t.Log(transFS.localizer)

	err = transFS.LoadTranslateFileFS(LocaleFS, "active.en_US.toml")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(transFS.tagMap)
	t.Log(transFS.currentTag)
	t.Log(transFS.localizer)

	err = transFS.SetCurrentTranslation("en_US")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(transFS.tagMap)
	t.Log(transFS.currentTag)
	t.Log(transFS.localizer)

}

func TestTranslator_Translate(t *testing.T) {
	transFS, err := NewTranslatorFS("active", LocaleFS, "active.zh_CN.toml")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	err = transFS.LoadTranslateFileFS(LocaleFS, "active.en_US.toml")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	t.Log(transFS.tagMap)
	t.Log(transFS.currentTag)
	t.Log(transFS.localizer)

	msg1 := transFS.Translate("ID_TEST", "fallback")
	t.Log("msg1:", msg1)

	err = transFS.SetCurrentTranslation("en_US")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log("\n")
	t.Log(transFS.tagMap)
	t.Log(transFS.currentTag)
	t.Log(transFS.localizer)
	msg2 := transFS.Translate("ID_TEST", "fallback")
	t.Log("msg2:", msg2)

}
