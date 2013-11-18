package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type generateFlag struct {
	Usage string
}

func (f generateFlag) String() string {
	return fmt.Sprintf("--generate, -g\t%v", f.Usage)
}

func (f generateFlag) Apply(set *flag.FlagSet) {
	set.Bool("g", false, f.Usage)
	set.Bool("generate", false, f.Usage)
}

func generateNewTask() (fileName string, err error) {
	sourceDir, err := os.Getwd()
	if err != nil {
		return
	}

	pkgName := filepath.Base(sourceDir)
	fileName = fmt.Sprintf("%s_task.go", pkgName)
  outfile := filepath.Join(sourceDir, fileName)
	f, err := os.Create(outfile)
	if err != nil {
		return
	}
	defer f.Close()

	w := exampleWriter{Pkg: pkgName}
	err = w.Write(f)

	return
}
