package main

import (
	"src/config"
	"src/flags"
	"src/pipe"
	"flag"
	"strings"
)

func main() {
	filename := flags.Filename
	flag.Parse()
	if len(*filename) == 0 {
		panic("filename not provided with -f flag")
	}
	ft := config.GetFT(*filename, config.SHELL)
	in, err := pipe.In()
	if err != nil {
		panic(err)
	}

	comment := ft.CmtStyle
	parts := strings.Split(strings.TrimSuffix(comment, " "), " ")
	mp := len(parts) > 1

	var scom string
	var ecom string
	if mp {
		if len(parts[0]) > 0 {
			scom = parts[0] + " "
		}
		if len(parts[1]) > 0 {
			scom = " " + parts[1]
		}
	}
	
	lines := strings.Split(string(in), "\n")
	nlines := []string{}
	for _, line := range lines {
		if len(line) == 0 {
			nlines = append(nlines, line)
			continue
		}

		if mp {
			hbegin := strings.Contains(line, scom)
			hend := strings.Contains(line, ecom)
			if hbegin && hend {
				nline := strings.Replace(line, scom, "", 1)
				nline = strings.Replace(nline, ecom, "", 1)
				nlines = append(nlines, nline)
				continue
			}
		}

		fch := 0
		for _, ch := range line {
			if ch == ' ' || ch == '\t' {
				fch++
				continue
			}
			break
		}

		fcom := fch + len(comment)
		if len(line) > fcom && line[fch:fcom] == comment {
			nline := strings.Replace(line, comment, "", 1)
			nlines = append(nlines, nline)
			continue
		}

		if mp {
			nline := line[:fch] + scom + line[fch:] + ecom
			nlines = append(nlines, nline)
			continue
		}
		nline := line[:fch] + comment + line[fch:]
		nlines = append(nlines, nline)
	}
	out := strings.Join(nlines, "\n")
	if out[len(out)-1] == '\n' {
		out = out[:len(out)-1]
	}
	pipe.Out(out)
}
