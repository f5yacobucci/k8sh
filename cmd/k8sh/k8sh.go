package main

import (
	"bufio"
	"bytes"
	"fmt"
	"k8sh/internal/executor"
	"k8sh/internal/parser"
	"k8sh/internal/symtab"
	"os"
)

func promptOne() {
	entry := symtab.GetSymbolEntry("PS1")
	if entry != nil {
		fmt.Fprintf(os.Stdout, "%s", entry.GetValue())
	} else {
		fmt.Fprintf(os.Stdout, "!!!$ ")
	}
}
func promptTwo() {
	entry := symtab.GetSymbolEntry("PS2")
	if entry != nil {
		fmt.Fprintf(os.Stdout, "%s", entry.GetValue())
	} else {
		fmt.Fprintf(os.Stdout, "!!!> ")
	}
}

func initShell() {
	entry := symtab.AddSymbol("CWD")
	entry.SetValue("/")

	entry = symtab.AddSymbol("PS1")
	entry.SetValue("$ ")

	entry = symtab.AddSymbol("PS2")
	entry.SetValue("> ")

	home, ok := os.LookupEnv("HOME")
	if !ok {
		fmt.Fprintf(os.Stderr, "MUST have $HOME set")
		os.Exit(1)
	}
	entry = symtab.AddSymbol("HOME")
	entry.SetValue(home)

	config, ok := os.LookupEnv("KUBECONFIG")
	if ok {
		entry = symtab.AddSymbol("KUBECONFIG")
		entry.SetValue(config)
	} else {
		config = home + "/.kube/config"
		entry = symtab.AddSymbol("KUBECONFIG")
		entry.SetValue(config)
	}
}

func main() {
	initShell()

	promptOne()

	var buf bytes.Buffer
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			promptOne()
			continue
		}
		if line[len(line)-1] == '\\' {
			buf.Write(line[:len(line)-1])
			promptTwo()
			continue
		}
		buf.Write(line)

		cmd := parser.Parse(buf)
		err := executor.DoSimpleCommand(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
		}

		buf.Reset()
		promptOne()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "IO Error: %v\n", err)
	}
}
