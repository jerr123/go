// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
)

// This program generates tests to verify that copying operations
// copy the data they are supposed to and clobber no adjacent values.

// run as `go run copyGen.go`.  A file called copy.go
// will be written into the parent directory containing the tests.

var sizes = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 16, 17, 23, 24, 25, 31, 32, 33, 63, 64, 65, 1023, 1024, 1025, 1024 + 7, 1024 + 8, 1024 + 9, 1024 + 15, 1024 + 16, 1024 + 17}

var usizes = [...]int{2, 3, 4, 5, 6, 7}

func main() {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "// Code generated by gen/copyGen.go. DO NOT EDIT.\n\n")
	fmt.Fprintf(w, "package main\n")
	fmt.Fprintf(w, "import \"testing\"\n")

	for _, s := range sizes {
		// type for test
		fmt.Fprintf(w, "type T%d struct {\n", s)
		fmt.Fprintf(w, "  pre [8]byte\n")
		fmt.Fprintf(w, "  mid [%d]byte\n", s)
		fmt.Fprintf(w, "  post [8]byte\n")
		fmt.Fprintf(w, "}\n")

		// function being tested
		fmt.Fprintf(w, "//go:noinline\n")
		fmt.Fprintf(w, "func t%dcopy_ssa(y, x *[%d]byte) {\n", s, s)
		fmt.Fprintf(w, "  *y = *x\n")
		fmt.Fprintf(w, "}\n")

		// testing harness
		fmt.Fprintf(w, "func testCopy%d(t *testing.T) {\n", s)
		fmt.Fprintf(w, "  a := T%d{[8]byte{201, 202, 203, 204, 205, 206, 207, 208},[%d]byte{", s, s)
		for i := 0; i < s; i++ {
			fmt.Fprintf(w, "%d,", i%100)
		}
		fmt.Fprintf(w, "},[8]byte{211, 212, 213, 214, 215, 216, 217, 218}}\n")
		fmt.Fprintf(w, "  x := [%d]byte{", s)
		for i := 0; i < s; i++ {
			fmt.Fprintf(w, "%d,", 100+i%100)
		}
		fmt.Fprintf(w, "}\n")
		fmt.Fprintf(w, "  t%dcopy_ssa(&a.mid, &x)\n", s)
		fmt.Fprintf(w, "  want := T%d{[8]byte{201, 202, 203, 204, 205, 206, 207, 208},[%d]byte{", s, s)
		for i := 0; i < s; i++ {
			fmt.Fprintf(w, "%d,", 100+i%100)
		}
		fmt.Fprintf(w, "},[8]byte{211, 212, 213, 214, 215, 216, 217, 218}}\n")
		fmt.Fprintf(w, "  if a != want {\n")
		fmt.Fprintf(w, "    t.Errorf(\"t%dcopy got=%%v, want %%v\\n\", a, want)\n", s)
		fmt.Fprintf(w, "  }\n")
		fmt.Fprintf(w, "}\n")
	}

	for _, s := range usizes {
		// function being tested
		fmt.Fprintf(w, "//go:noinline\n")
		fmt.Fprintf(w, "func tu%dcopy_ssa(docopy bool, data [%d]byte, x *[%d]byte) {\n", s, s, s)
		fmt.Fprintf(w, "  if docopy {\n")
		fmt.Fprintf(w, "    *x = data\n")
		fmt.Fprintf(w, "  }\n")
		fmt.Fprintf(w, "}\n")

		// testing harness
		fmt.Fprintf(w, "func testUnalignedCopy%d(t *testing.T) {\n", s)
		fmt.Fprintf(w, "  var a [%d]byte\n", s)
		fmt.Fprintf(w, "  t%d := [%d]byte{", s, s)
		for i := 0; i < s; i++ {
			fmt.Fprintf(w, " %d,", s+i)
		}
		fmt.Fprintf(w, "}\n")
		fmt.Fprintf(w, "  tu%dcopy_ssa(true, t%d, &a)\n", s, s)
		fmt.Fprintf(w, "  want%d := [%d]byte{", s, s)
		for i := 0; i < s; i++ {
			fmt.Fprintf(w, " %d,", s+i)
		}
		fmt.Fprintf(w, "}\n")
		fmt.Fprintf(w, "  if a != want%d {\n", s)
		fmt.Fprintf(w, "    t.Errorf(\"tu%dcopy got=%%v, want %%v\\n\", a, want%d)\n", s, s)
		fmt.Fprintf(w, "  }\n")
		fmt.Fprintf(w, "}\n")
	}

	// boilerplate at end
	fmt.Fprintf(w, "func TestCopy(t *testing.T) {\n")
	for _, s := range sizes {
		fmt.Fprintf(w, "  testCopy%d(t)\n", s)
	}
	for _, s := range usizes {
		fmt.Fprintf(w, "  testUnalignedCopy%d(t)\n", s)
	}
	fmt.Fprintf(w, "}\n")

	// gofmt result
	b := w.Bytes()
	src, err := format.Source(b)
	if err != nil {
		fmt.Printf("%s\n", b)
		panic(err)
	}

	// write to file
	err = ioutil.WriteFile("../copy_test.go", src, 0666)
	if err != nil {
		log.Fatalf("can't write output: %v\n", err)
	}
}
