package symtab

import (
	"fmt"
	"os"
	"reflect"
)

type SymbolType int

const (
	SymString SymbolType = iota
	SymFunc              = iota
)

type SymbolEntry struct {
	name  string
	kind  SymbolType
	val   reflect.Value
	flags int
	// README.md: parser.Node* funcBody
}

type SymbolMap map[string]*SymbolEntry
type SymbolTable struct {
	level   int
	symbols SymbolMap
}
type SymbolStack struct {
	stack  []*SymbolTable
	global *SymbolTable
	local  *SymbolTable
}

var tables *SymbolStack
var tableLevel int

func init() {
	initTable := &SymbolTable{
		0,
		make(SymbolMap),
	}
	tables = &SymbolStack{
		global: initTable,
		local:  initTable,
		stack:  []*SymbolTable{initTable},
	}
}

func NewTable(level int) *SymbolTable {
	table := &SymbolTable{
		level:   level,
		symbols: make(SymbolMap),
	}
	return table
}

func DumpLocalTable() {
	locals := tables.local

	indent := locals.level * 4
	fmt.Fprintf(os.Stderr, "%*sSymbol table [Level %d]:\n", indent, " ", locals.level)
	fmt.Fprintf(os.Stderr, "%*s===========================\n", indent, " ")
	fmt.Fprintf(os.Stderr, "%*s            Symbol                    Val\r\n", indent, " ")
	fmt.Fprintf(os.Stderr, "%*s-------------------------------- ------------\r\n", indent, " ")

	for k, v := range locals.symbols {
		fmt.Fprintf(os.Stderr, "%*s %-32s '%s'\n", indent, " ",
			k, v.val.String())
	}
	fmt.Fprintf(os.Stderr, "%*s-------------------------------- ------------\r\n", indent, " ")
}

func AddSymbol(symbol string) *SymbolEntry {
	if len(symbol) == 0 {
		return nil
	}

	locals := tables.local
	if v, ok := locals.symbols[symbol]; ok {
		return v
	}

	entry := SymbolEntry{
		name: symbol,
	}

	locals.symbols[symbol] = &entry
	return &entry
}

func RemoveSymbol(symbol string, table SymbolMap) {
	delete(table, symbol)
}

func GetSymbolEntry(symbol string) *SymbolEntry {
	for i := len(tables.stack) - 1; i >= 0; i-- {
		if v, ok := tables.stack[i].symbols[symbol]; ok {
			return v
		}
	}

	return nil
}

func (e *SymbolEntry) GetValue() string {
	return e.val.String()
}
func (e *SymbolEntry) SetValue(val string) {
	e.val = reflect.ValueOf(val)
}

// Stack Functions
func StackAdd(table *SymbolTable) {
	tables.stack = append(tables.stack, table)
	tables.local = table
}

func StackPush() *SymbolTable {
	tableLevel += 1
	table := NewTable(tableLevel)
	StackAdd(table)
	return table
}

func StackPop() *SymbolTable {
	if len(tables.stack) == 0 {
		return nil
	}

	table := tables.stack[len(tables.stack)-1]
	tableLevel -= 1
	tables.stack = tables.stack[:len(tables.stack)-1]

	if len(tables.stack) == 0 {
		tables.local = nil
		tables.global = nil
	} else {
		tables.local = tables.stack[len(tables.stack)-1]
	}
	return table
}

func GetLocals() *SymbolTable {
	return tables.local
}

func GetGlobals() *SymbolTable {
	return tables.global
}

func GetStack() *SymbolStack {
	return tables
}
