package docit

import (
	"unicode"
)

type Path string

func (p Path) FileName() string {
	l := len(p)
	name := make([]byte, 0)

	fileExtensionPoint := 0
	for i := l - 1; i > 0; i-- {
		if p[i] == '.' {
			fileExtensionPoint = i
			break
		}
	}

	for i := fileExtensionPoint - 1; i > 0; i-- {
		if p[i] == '/' {
			break
		}

		name = append(name, p[i])
	}

	nameLength := len(name)
	for i := 0; i < nameLength/2; i++ {
		tmp := name[i]
		name[i] = name[nameLength-i-1]
		name[nameLength-i-1] = tmp
	}

	firstSymbol := rune(name[0])
	if !(unicode.IsLetter(firstSymbol) || unicode.IsNumber(firstSymbol)) {
		return string(name[1:])
	} else {
		return string(name)
	}
}
