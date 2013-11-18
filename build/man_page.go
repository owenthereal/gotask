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
	sections, err := p.readSections()
	if err != nil {
		return
	}

	mp = &manPage{}
	if nameAndUsage, ok := sections["NAME"]; ok {
		mp.Name, mp.Usage = p.splitNameAndUsage(nameAndUsage)
	}
	mp.Description = sections["DESCRIPTION"]

	return
}

func (p *manPageParser) readSections() (sections map[string]string, err error) {
	sections = make(map[string]string)
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
					sections[heading] = concatHeadingContent(content)
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
		sections[heading] = concatHeadingContent(content)
	}

	if err == io.EOF {
		err = nil
	}

	return
}

func (p *manPageParser) splitNameAndUsage(nameAndUsage string) (name, usage string) {
	s := strings.SplitN(nameAndUsage, " - ", 2)
	if len(s) == 1 {
		name = ""
		usage = strings.TrimSpace(s[0])
	} else {
		name = strings.TrimSpace(s[0])
		usage = strings.TrimSpace(s[1])
	}

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
