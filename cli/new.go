package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type newFlag struct {
	Usage string
}

func (f newFlag) String() string {
	return fmt.Sprintf("--new, -n\t%v", f.Usage)
}

func (f newFlag) Apply(set *flag.FlagSet) {
	set.Bool("n", false, f.Usage)
	set.Bool("new", false, f.Usage)
}

func generateNewTask() (err error) {
	sourceDir, err := os.Getwd()
	if err != nil {
		return
	}

	pkgName := filepath.Base(sourceDir)
	fileName := fmt.Sprintf("%s_task.go", pkgName)
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
