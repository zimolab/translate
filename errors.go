package translate

import "fmt"

type IllegalFilename struct {
	filename string
}

func (e IllegalFilename) Error() string {
	return fmt.Sprintf("illegal translate filename: %s", e.filename)
}

func illegalFilename(filename string) IllegalFilename {
	return IllegalFilename{filename: filename}
}

type TranslationNotLoaded struct {
	name string
}

func (e TranslationNotLoaded) Error() string {
	return fmt.Sprintf("translation loaded: %s", e.name)
}

func notLoaded(name string) TranslationNotLoaded {
	return TranslationNotLoaded{name: name}
}
