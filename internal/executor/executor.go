package executor

import (
	"errors"
	"k8sh/internal/builtin"
	"k8sh/internal/parser"
	"k8sh/internal/symtab"
)

func DoSimpleCommand(cmd *parser.Node) error {
	if cmd == nil {
		return errors.New("cmd node is nil")
	}
	if cmd.Kind != parser.NodeCommand {
		return errors.New("cmd node is not a command")
	}
	if cmd.FirstChild == nil {
		return errors.New("cmd has no children")
	}

	var argc int
	var argv []string
	for child := cmd.FirstChild; child != nil; child = child.NextSibling {
		str := child.Val.String()

		str = tryExpand(str)

		argv = append(argv, str)
	}
	argc = len(argv)

	// check builtins
	if v, ok := builtin.BuiltIns[argv[0]]; ok {
		return v.Exec(argc, argv)
	}

	// plugins unsupported
	return errors.New("command not found")
}

func tryExpand(str string) string {
	for index, character := range str {
		switch character {
		case '$':
			// README.md handle quotes and word expansion
			maybeVar := str[index+1 : len(str)]
			entry := symtab.GetSymbolEntry(maybeVar)
			if entry != nil {
				return entry.GetValue()
			}
		}
	}

	return str
}
