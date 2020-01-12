package domain

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"
)

func TestTree_Has(t *testing.T) {
	domainTree := DefaultTree()

	domainFile, err := os.Open("tmp/domain_test.txt")

	if err != nil {
		t.Fatal(err)
	}

	buf := bufio.NewReader(domainFile)

	for {

		line, err := buf.ReadString('\n')

		if err != nil && err != io.EOF {
			continue
		}

		line = strings.TrimSpace(line)

		domainTree.Insert(line)
		if err == io.EOF {
			break
		}
	}

	if ! domainTree.Has("a.com") {
		t.FailNow()
	}
	if ! domainTree.Has("b.com") {
		t.FailNow()
	}
	if ! domainTree.Has("a.b.com") {
		t.FailNow()
	}
	if domainTree.Has("c.com") {
		t.FailNow()
	}
	if ! domainTree.Has("a.c.com") {
		t.FailNow()
	}
	if ! domainTree.Has("a.cn") {
		t.FailNow()
	}
}
