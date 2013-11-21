package build

import (
	"bufio"
	"bytes"
	"github.com/jingweno/gotask/task"
	"io"
	"regexp"
	"strings"
)

type manPage struct {
	Name        string
	Usage       string
	Description string
	Flags       []task.Flag
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
	if opts, ok := sections["OPTIONS"]; ok {
		mp.Flags, err = p.parseOptions(opts)
		if err != nil {
			return
		}
	}
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
					sections[heading] = concatSectionContent(content)
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
		sections[heading] = concatSectionContent(content)
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

func (p *manPageParser) parseOptions(optsStr string) (flags []task.Flag, err error) {
	reader := bufio.NewReader(bytes.NewReader([]byte(optsStr)))
	flagRegexp := regexp.MustCompile(`\-?\-(\w+),?`)

	var (
		line, name string
		content    []string
	)
	for err == nil {
		line, err = readLine(reader)
		if flagRegexp.MatchString(line) {
			if name != line {
				if name != "" {
					flags = append(flags, task.NewBoolFlag(name, concatFlagContent(content)))
				}

				var fstrs []string
				for _, fstr := range flagRegexp.FindAllStringSubmatch(line, -1) {
					fstrs = append(fstrs, fstr[1])
				}

				name = strings.Join(fstrs, ", ")
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
	if name != "" {
		flags = append(flags, task.NewBoolFlag(name, concatFlagContent(content)))
	}

	if err == io.EOF {
		err = nil
	}

	return
}

func concatSectionContent(content []string) string {
	return strings.TrimSpace(strings.Join(content, "\n   "))
}

func concatFlagContent(content []string) string {
	return strings.TrimSpace(strings.Join(content, "\n"))
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
