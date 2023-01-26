package translate

import "fmt"

type ErrIllegalTranslateFilename struct {
	filename string
}

func (e ErrIllegalTranslateFilename) Error() string {
	return fmt.Sprintf("illegal translate filename: %s", e.filename)
}

func illegalFilename(filename string) ErrIllegalTranslateFilename {
	return ErrIllegalTranslateFilename{filename: filename}
}

type ErrTranslationNotLoaded struct {
	name string
}

func (e ErrTranslationNotLoaded) Error() string {
	return fmt.Sprintf("translation loaded: %s", e.name)
}

func translationNotLoaded(name string) ErrTranslationNotLoaded {
	return ErrTranslationNotLoaded{name: name}
}

type ErrIllegalLanguageTag struct {
	tag string
}

func (e ErrIllegalLanguageTag) Error() string {
	return fmt.Sprintf("illegal language tag: %s", e.tag)
}

func illegalLanguageTag(tag string) ErrIllegalLanguageTag {
	return ErrIllegalLanguageTag{tag: tag}
}
