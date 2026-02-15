package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Vi style motion emulation.
//
//	count command
//
// Where
//
//	count := \d*
//	command := [hjklgGqfx]+
//
// Enables behaviour like 25j for moving 25 mails down, 2gf for going to the
// second attachment and gx for opening an url or file path via the operating
// systems open behaviour
type vi struct {
	modifier uint
	command  strings.Builder
}

var validCommands = map[string]struct{}{
	"h":  {},
	"j":  {},
	"k":  {},
	"l":  {},
	"G":  {},
	"gg": {},
	"gf": {},
	"gx": {},
	"q":  {},
}

func isPrefix(s string) bool {
	for cmd := range validCommands {
		if strings.HasPrefix(cmd, s) {
			return true
		}
	}
	return false
}

// represent a fully detected vi motion
type viMessage struct {
	modifier uint
	command  string
}

func (v *vi) reset() {
	v.modifier = 0
	v.command.Reset()
}

func (v *vi) pending() string {
	if v.modifier == 0 || v.modifier == 1 {
		return fmt.Sprint(v.command.String())
	} else {
		return fmt.Sprint(v.modifier, v.command.String())
	}
}

// convert the current vi state into a viMessage model.Update can deal with
func (v *vi) toViMessage() viMessage {
	if v.modifier < 1 {
		v.modifier = 1
	}
	msg := viMessage{
		modifier: v.modifier,
		command:  v.command.String(),
	}
	return msg
}

// update the vi state
func (v *vi) update(msg tea.KeyMsg) (viMessage, bool) {
	switch msg.Type {
	case tea.KeyEsc:
		v.reset()
	case tea.KeyRunes:
		if len(msg.Runes) != 1 {
			return viMessage{}, false
		}
		k := msg.Runes[0]
		switch {
		case k >= '0' && k <= '9':
			v.modifier = v.modifier*10 + uint(k-'0')
		default:
			switch k {
			case 'k', 'j', 'h', 'l', 'G', 'q', 'f', 'g', 'x':
				v.command.WriteRune(k)

				cmd := v.command.String()
				if _, ok := validCommands[cmd]; ok {
					vimsg := v.toViMessage()
					v.reset()
					return vimsg, true
				}

				if !isPrefix(cmd) {
					v.reset()
				}
			}
		}
	}

	return viMessage{}, false
}
