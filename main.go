/*
Copyright (c) 2013, Aulus Egnatius Varialus <varialus@gmail.com>

Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
*/

/*
c2go is based on Go2C at https://github.com/xyproto/c2go
Copyright (c) 2011-2013, Alexander Rødseth
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name of the author nor the
      names of its contributors may be used to endorse or promote products
      derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

/*
Go2C is based on the c-to-c.py example from pycparser by Eli Bendersky, and uses pycparser extensively.
Copyright (C) 2008-2011, Eli Bendersky
License: BSD
*/

// https://github.com/xyproto/c2go/blob/master/c2go.py

package main

import (
	"github.com/varialus/bsd/temporary_translation_utilities"
)

func main() {

	REPLACEMENT_FUNCTIONS := map[string]string{
		"stat": "syscall.Stat",
		"access": "syscall.Access",
		"rand": "rand.Float64",
	}

	REPLACEMENT_TYPES := map[string]string{
		"static struct stat": "syscall.Stat_t",
		"struct timeval": "syscall.Timeval",
		"char *": "CString",
		"char": "byte",
		"unsigned char": "byte",
		"int *": "*int",
		"int": "int",
		"unsigned int": "uint",
		"unsigned int *": "*uint",
		"void": "",
		"short": "int16",
		"short *": "*int16",
		"unsigned short": "uint16",
		"unsigned short *": "*uint16",
		"float": "float32",
		"float *": "*float32",
		"double": "float64",
		// TODO: check if int is int32 in Go
		"long": "int",
		// TODO: Needs a better plan for static
		"static long": "int",
		"static int": "int",
	}

	REPLACEMENT_MACROS := map[string][3]string{
		"S_ISDIR" : [3]string{"syscall", "(((", ") & syscall.S_IFMT) == syscall.S_IFDIR)"},
	}

	REPLACEMENT_DEFS := map[string]int{
		"F_OK" : 0,
		"X_OK" : 1,
		"W_OK" : 2,
		"R_OK" : 4,
	}

	CUSTOM_FUNCTIONS := map[string][2]string{
		"atoi":   [2]string{"strconv",  "func atoi(a string) int {\n\tv, _ := strconv.Atoi(a)\n\treturn v\n}"},
		"sleep":  [2]string{"time",     "func sleep(sec int64) {\n\ttime.Sleep(1e9 * sec)\n}"},
		"getchar":[2]string{"fmt",      "func getchar() byte {\n\tvar b byte\n\tfmt.Scanf(\"%c\", &b)\n\treturn b\n}"},
		"putchar":[2]string{"fmt",      "func putchar(b byte) {\n\tfmt.Printf(\"%c\", b)\n}"},
		"abs":    [2]string{"",         "func abs(a int) int {\n\tif a >= 0 {\n\t\treturn a\n\t}\n\treturn -a\n}"},
		"strcpy": [2]string{"CString",  "func strcpy(a *CString, b CString) {\n\t*a = b\n}"},
		"strcmp": [2]string{"CString",  "func strcmp(acs, bcs CString) int {\n\ta := acs.ToString()\n\tb := bcs.ToString()\n\tif a == b {\n\t\treturn 0\n\t}\n\talen := len(a)\n\tblen := len(b)\n\tminlen := blen\n\tif alen < minlen {\n\t\tminlen = alen\n\t}\n\tfor i := 0; i < minlen; i++ {\n\t\tif a[i] > b[i] {\n\t\t\treturn 1\n\t\t} else if a[i] < b[i] {\n\t\t\treturn -1\n\t\t}\n\t}\n\tif alen > blen {\n\t\treturn 1\n\t}\n\treturn -1\n}"},
		"strlen": [2]string{"CString",  "func strlen(c CString) int {\n\treturn len(c.ToString())\n}"},
		"printf": [2]string{"fmt",      "func printf(format CString, a ...interface{}) {\n\tfmt.Printf(format.ToString(), a...)\n}"},
		"scanf":  [2]string{"fmt",      "func scanf(format CString, a ...interface{}) {\n\tfmt.Scanf(format.ToString(), a...)\n}"},
		"b2i":    [2]string{"",         "func b2i(b bool) int {\n\tif b{\n\t\treturn 1\n\t}\n\treturn 0\n}"},
	}

	SKIP_INCLUDES := [1]string{"CString"}

	WHOLE_PROGRAM_REPLACEMENTS := map[string]string{
		"'fmt.Printf(\"\n\")": "fmt.Println()",
		"func main(argc int, argv *[]CString) int {": "func main() {\n\tflag.Parse()\n\targv := flag.Args()\n\targc := len(argv)+1\n",
		"argv[": "argv[-1+",
		"for (1)": "for",
		"\n\n\n": "\n",
		// TODO: Find a sensible way to figure out when a program wants strings and when it wants byte arrays
		"*[]byte": "CString",
		//"*[128]byte = new([128]byte)": "*string",
		//
		"'\\0'": "'\\x00'",
		"int = 0\n": "int\n",
	}

	temporary_translation_utilities.Use_vars_so_compiler_does_not_complain(
		REPLACEMENT_FUNCTIONS,
		REPLACEMENT_TYPES,
		REPLACEMENT_MACROS,
		REPLACEMENT_DEFS,
		CUSTOM_FUNCTIONS,
		SKIP_INCLUDES,
		WHOLE_PROGRAM_REPLACEMENTS)
}
