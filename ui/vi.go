package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// `<modifier><operator>` enables behaviour like 25j for moving 25 mails down
// or 2gf for going to the second attachemnt
type vi struct {
	modifier uint
	operator rune
}

// represent an vi like modal operation
type viMessage struct {
	modifier uint
	operator rune
}

func (v *vi) reset() {
	v.modifier = 0
	v.operator = 0
}

func (v *vi) pending() string {
	if v.modifier == 0 || v.modifier == 1 {
		return fmt.Sprint(string(v.operator))
	} else {
		return fmt.Sprint(v.modifier, string(v.operator))
	}
}

// convert the current vi state into a viMessage model.Update can deal with
func (v *vi) toViMessage() viMessage {
	if v.modifier < 1 {
		v.modifier = 1
	}
	msg := viMessage{
		modifier: v.modifier,
		operator: v.operator,
	}
	return msg
}

// update the vi state
func (v *vi) update(msg tea.KeyMsg) (viMessage, bool) {
	switch msg.Type {
	case tea.KeyRunes:
		if len(msg.Runes) != 1 {
			return viMessage{}, false
		}
		k := msg.Runes[0]
		switch {
		case k >= '0' && k <= '9':
			v.modifier = v.modifier*10 + uint(k-'0')
		case k == 'k',
			k == 'j',
			k == 'h',
			k == 'l',
			k == 'G',
			k == 'q':
			v.operator = k
			return v.toViMessage(), true
		}
	}

	return viMessage{}, false
}
