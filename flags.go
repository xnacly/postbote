package postbote

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Flags struct {
	Verbose bool
	Account string
	Folder  string
	Command string
	Config  string
}

func (c *Flags) FromArgs() {
	flag.BoolVar(&c.Verbose, "verbose", false, "enable verbose logging")
	flag.StringVar(&c.Account, "account", "", "which account to connect to, if multiple exist, otherwise the singular account is selected")
	flag.StringVar(&c.Folder, "dir", "INBOX", "which mail dir to start in")
	flag.StringVar(&c.Command, "cmd", "", "start postbote, execute :<cmd>, print output to stdout and exit")
	flag.StringVar(&c.Config, "config", "~/.config/postbote.toml", "change the path to the postbote configuration file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Postbote A minimalist opinionated mutt alternative Mail client inspired by vi, built for the terminal

Usage:
  %s [options]

Options:
`, filepath.Base(os.Args[0]))

		flag.PrintDefaults()
	}
	flag.Parse()
}
