package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/xnacly/postbote"
	"github.com/xnacly/postbote/ui"
)

func main() {
	f := postbote.Flags{}
	f.FromArgs()

	log.SetOutput(os.Stderr)
	if f.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("Parsed flags", "flags", f)

	if err := ui.Run(f); err != nil {
		panic(err)
	}
}
