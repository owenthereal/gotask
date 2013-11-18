package cli

import (
	"flag"
	"fmt"
	"github.com/jingweno/gotask/build"
	"os"
	"path/filepath"
)

type compileFlag struct {
	Usage string
}

func (f compileFlag) String() string {
	return fmt.Sprintf("--compile, -c\t%v", f.Usage)
}

func (f compileFlag) Apply(set *flag.FlagSet) {
	set.Bool("c", false, f.Usage)
	set.Bool("compile", false, f.Usage)
}

func compileTasks(isDebug bool) (err error) {
	sourceDir, err := os.Getwd()
	if err != nil {
		return
	}

	fileName := fmt.Sprintf("%s.task", filepath.Base(sourceDir))
	outfile := filepath.Join(sourceDir, fileName)

	err = build.Compile(sourceDir, outfile, isDebug)
	return
}
