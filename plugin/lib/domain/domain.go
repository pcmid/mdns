package domain

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Domain string

type Tree struct {
	mark uint8
	sub  domainMap
}

type domainMap map[Domain]*Tree

func (d Domain) nextLevel() Domain {
	if pointIndex := strings.LastIndex(string(d), "."); pointIndex == -1 {
		return ""
	} else {
		return d[:pointIndex]
	}
}

func (d Domain) topLevel() Domain {
	if pointIndex := strings.LastIndex(string(d), "."); pointIndex == -1 {
		return d
	} else {
		return d[pointIndex+1:]
	}
}

func DefaultTree() *Tree {
	return NewDomainTree()
}

func NewDomainTree() (dt *Tree) {
	dt = new(Tree)
	dt.sub = make(domainMap)
	return
}

func (dt *Tree) Has(d Domain) bool {

	if strings.LastIndexByte(string(d), '.') == len(d)-1 {
		d = Domain(strings.TrimRight(string(d), "."))
	}

	if len(dt.sub) == 0 {
		return false
	}

	return dt.has(d)
}

func (dt *Tree) has(d Domain) bool {

	if len(dt.sub) == 0 {
		return true
	}

	if sub, ok := dt.sub[d.topLevel()]; ok {
		return sub.has(d.nextLevel())
	}
	return false
}

func (dt *Tree) insert(sections []Domain) {

	if len(sections) == 0 {
		return
	}

	lastIndex, lastSec := len(sections)-1, sections[len(sections)-1]
	if lastSec == "" {
		return
	}

	if sec, ok := dt.sub[lastSec]; ok {
		sec.insert(sections[:lastIndex])
	} else {
		dt.sub[lastSec] = NewDomainTree()
		dt.sub[lastSec].insert(sections[:lastIndex])
	}
}

func (dt *Tree) Insert(d string) {
	sections := strings.Split(d, ".")
	if len(sections) == 0 {
		return
	}

	domainSec := make([]Domain, len(sections))

	for i := range sections {
		domainSec[i] = Domain(sections[i])
	}

	dt.insert(domainSec)
}

func TreeFromFile(file string) (dt *Tree, err error) {
	dt = DefaultTree()
	domainFile, errR := os.Open(file)
	defer func() {
		_ = domainFile.Close()
	}()

	if errR != nil {
		return dt, errR
	}

	buf := bufio.NewReader(domainFile)

	for {

		line, err := buf.ReadBytes('\n')

		if err != nil && err != io.EOF {
			continue
		}

		lineB := make([]byte, len(line))
		copy(lineB, line)

		lineS := strings.TrimSpace(string(lineB))

		dt.Insert(lineS)

		if err == io.EOF {
			break
		}
	}
	buf = nil
	return

}
