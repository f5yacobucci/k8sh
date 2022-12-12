package parser

import (
	"bufio"
	"bytes"
)

func Parse(buf bytes.Buffer) *Node {
	if buf.Len() == 0 {
		return nil
	}

	cmd := NewNode(NodeCommand)

	scanner := bufio.NewScanner(&buf)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := NewNode(NodeVar)
		text := scanner.Text()
		SetNodeValueStr(word, text)
		AddChild(cmd, word)
	}
	return cmd
}
