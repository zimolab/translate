package translate

type IllegalFilename struct {
	filename string
}

func (e IllegalFilename) Error() string {
	return f("illegal filename: %s", e.filename)
}

func illegalFilename(filename string) IllegalFilename {
	return IllegalFilename{filename: filename}
}

type NotLoaded struct {
	name string
}

func (e NotLoaded) Error() string {
	return f("not loaded: %s", e.name)
}

func notLoaded(name string) NotLoaded {
	return NotLoaded{name: name}
}

type MissingDisplayName struct {
	path string
}

func (e MissingDisplayName) Error() string {
	return f("missing display name: %s", e.path)
}

func missingDisplayName(path string) MissingDisplayName {
	return MissingDisplayName{path: path}
}

type LocaleNotSet struct {
	msg string
}

func (e LocaleNotSet) Error() string {
	return e.msg
}

func localeNotSet() LocaleNotSet {
	return LocaleNotSet{
		msg: "locale not set",
	}
}
