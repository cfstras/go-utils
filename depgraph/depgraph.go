// This is a simple command line tool to create a dependency graph from a path.
// Usage:
//  depgraph <root package> | dot -Tsvg > graph.svg
//
// Patches welcome.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

var stdLib []string = getStdLib()

type goList struct {
	Imports []string
}

var done = make(map[string]bool)

func getDeps(p string) []string {
	o, err := exec.Command("go", "list", "-json", p).Output()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
	}

	list := goList{}
	json.Unmarshal(o, &list)

	return list.Imports
}

func printRecursive(p string) {
	done[p] = true

	for _, d := range getDeps(p) {
		if containsString(stdLib, d) {
			continue
		}
		fmt.Printf("\t\"%s\" -- \"%s\";\n", p, d)

		if !done[d] {
			printRecursive(d)
		}
	}
}

func main() {
	fmt.Println("graph G {")
	printRecursive(os.Args[1])
	fmt.Println("}")
}

// this is an ugly hack, I know.
// TODO make this prettier
func getStdLib() []string {
	return []string{
		"bufio",
		"bytes",
		"crypto",
		"crypto/aes",
		"crypto/cipher",
		"crypto/md5",
		"crypto/rannd",
		"crypto/sha1",
		"database/sql",
		"database/sql/driver",
		"encoding/binary",
		"encoding/json",
		"errors",
		"flag",
		"fmt",
		"html",
		"html/template",
		"io",
		"io/ioutil",
		"log",
		"math",
		"math/big",
		"net",
		"net/http",
		"os",
		"os/signal",
		"os/user",
		"path",
		"path/filepath",
		"reflect",
		"regexp/syntax",
		"runtime",
		"sort",
		"strconv",
		"strings",
		"sync",
		"syscall",
		"testing",
		"time",
		"unicode",
		"unsafe"}
}

func containsString(haystack []string, needle string) bool {
	for _, hay := range haystack {
		if hay == needle {
			return true
		}
	}
	return false
}
