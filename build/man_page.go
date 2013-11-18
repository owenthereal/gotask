package build

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
)

type manPage struct {
	Name        string
	Usage       string
	Description string
}

type manPageParser struct {
	Doc string
}

func (p *manPageParser) Parse() (mp *manPage, err error) {
	mp = &manPage{}
	result := make(map[string]string)
	headingRegexp := regexp.MustCompile(`^([A-Z]+)$`)
	reader := bufio.NewReader(bytes.NewReader([]byte(p.Doc)))

	var (
		line    string
		heading string
		content []string
	)
	for err == nil {
		line, err = readLine(reader)

		if headingRegexp.MatchString(line) {
			if heading != line {
				if heading != "" {
					result[heading] = concatHeadingContent(content)
				}

				heading = line
				content = []string{}
			}
		} else {
			if line != "" {
				line = strings.TrimSpace(line)
			}
			content = append(content, line)
		}
	}
	// the last one
	if heading != "" {
		result[heading] = concatHeadingContent(content)
	}

	if err == io.EOF {
		err = nil
	}

	// set NAME and USAGE
	if name, ok := result["NAME"]; ok {
		s := strings.SplitN(name, " - ", 2)
		if len(s) == 1 {
			mp.Name = ""
			mp.Usage = strings.TrimSpace(s[0])
		} else {
			mp.Name = strings.TrimSpace(s[0])
			mp.Usage = strings.TrimSpace(s[1])
		}
	}
	mp.Description = result["DESCRIPTION"]

	return
}

func concatHeadingContent(content []string) string {
	return strings.TrimSpace(strings.Join(content, "\n   "))
}

func readLine(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return string(ln), err
}
