package path

import (
	"errors"
	"k8sh/internal/symtab"
	"path"
	"strings"
)

type PathFlags int

const (
	Root      PathFlags = iota
	Namespace           = iota
	Resource            = iota
	Object              = iota
)

type Element struct {
	level PathFlags
	value string
	next  *Element
}

func IsAbs(p string) bool {
	return path.IsAbs(p)
}

func MakeAbsolute(p string) string {
	if IsAbs(p) {
		return p
	}

	// errcheck
	entry := symtab.GetSymbolEntry("CWD")
	cwd := entry.GetValue()

	ret := path.Join(cwd, p)

	return ret
}

func makeAbsoluteElements(p string) (*Element, error) {
	if !IsAbs(p) {
		return nil, errors.New("must be absolute path")
	}
	cleaned := path.Clean(p)

	var root *Element = &Element{
		level: Root,
	}
	var cur *Element
	for i, v := range strings.Split(cleaned, "/") {
		if v == "" {
			continue
		}
		e := &Element{
			level: PathFlags(i),
			value: v,
		}
		cur.next = e
		cur = e
	}
	return root, nil
}

func CheckPath(p string) (bool, *Element, error) {
	root, err := makeAbsoluteElements(p)
	if err != nil {
		return false, root, err
	}

	// TODO: Finish this
	return false, nil, nil
}
