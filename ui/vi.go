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

var validCommandList = []string{
	"h",
	"j",
	"k",
	"l",
	"G",
	"gg",
	"gf",
	"gx",
	"q",
	"/",
	"a",
}
var validCommands = map[string]struct{}{}
var validFirstRunes = map[rune]struct{}{}
var validPrefixes = map[string]struct{}{}

func init() {
	for _, cmd := range validCommandList {
		validCommands[cmd] = struct{}{}
		validFirstRunes[rune(cmd[0])] = struct{}{}
		// this is currently the best way i can think of, besides a trie, which
		// i dont think is necessary for a whole 11 commands
		for i := 1; i <= len(cmd); i++ {
			validPrefixes[cmd[:i]] = struct{}{}
		}
	}
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
		case k >= '0' && k <= '9' && v.command.Len() == 0:
			v.modifier = v.modifier*10 + uint(k-'0')
		default:
			if _, ok := validFirstRunes[k]; ok {
				v.command.WriteRune(k)
				cmd := v.command.String()
				if _, ok := validCommands[cmd]; ok {
					vimsg := v.toViMessage()
					v.reset()
					return vimsg, true
				}

				if _, ok := validPrefixes[cmd]; !ok {
					v.reset()
				}
			}
		}
	}

	return viMessage{}, false
}
