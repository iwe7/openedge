package module

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Flags command-line flags
type Flags struct {
	WorkDir string
	Config  string
	Help    bool
}

// ParseFlags parses the command-line flags
func ParseFlags(defaultConfigPath string) (*Flags, error) {
	f := new(Flags)
	cwd, err := os.Executable()
	if err != nil {
		return nil, err
	}
	cwd, err = filepath.EvalSymlinks(cwd)
	if err != nil {
		return nil, err
	}
	f.WorkDir = filepath.Dir(filepath.Dir(cwd))
	if defaultConfigPath == "" {
		f.Config = filepath.Join("etc", "openedge", "module.yml")
	} else {
		f.Config = defaultConfigPath
	}
	flag.StringVar(
		&f.WorkDir,
		"w",
		f.WorkDir,
		"working directory",
	)
	flag.StringVar(
		&f.Config,
		"c",
		f.Config,
		"config file path",
	)
	flag.BoolVar(
		&f.Help,
		"h",
		false,
		"show this help",
	)
	flag.Parse()
	f.WorkDir, err = filepath.Abs(f.WorkDir)
	if err != nil {
		return nil, err
	}
	err = os.Chdir(f.WorkDir)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// PrintUsage prints usage
func PrintUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Version of %s: %s\n", os.Args[0], Version)
	flag.Usage()
}
