package main

import (
	"bufio"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"unicode/utf8"
)

type LenString struct {
	n int
	s string
}

func main() {
	var longestUtf8 LenString
	var longestBytes LenString
	etlds := make(map[byte][]string)

	// web request
	h, err := http.Get("https://publicsuffix.org/list/public_suffix_list.dat")
	if err != nil {
		panic(err)
	}

	// scan lines
	sc := bufio.NewScanner(h.Body)
	for sc.Scan() {
		t := sc.Text()
		if t == "" || strings.HasPrefix(t, "//") {
			continue
		}
		l := utf8.RuneCountInString(t)
		if l > longestUtf8.n {
			longestUtf8 = LenString{l, t}
		}
		l = len(t)
		if l > longestBytes.n {
			longestBytes = LenString{l, t}
		}

		n := byte(strings.Count(t, ".")) + 1
		etlds[n] = append(etlds[n], t)
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("Longest UTF-8: %v\n", longestUtf8)
	fmt.Printf("Longest Bytes: %v\n", longestBytes)

	for i := 0; i < 10; i++ {
		if v, ok := etlds[byte(i)]; ok {
			sort.Slice(v, func(i, j int) bool {
				return utf8.RuneCountInString(v[i]) > utf8.RuneCountInString(v[j])
			})
			fmt.Printf("=== %d sections (%d) ===\n", i, len(v))
			l := len(v)
			if l > 5 {
				l = 5
			}
			for i := 0; i < l; i++ {
				fmt.Println(v[i])
			}
		}
	}
}
