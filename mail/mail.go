package mail

import "time"

type Folder struct {
	Name     string
	Path     string
	Messages []Message
}

type Message struct {
	UID         uint32
	Subject     string
	From        string
	Date        time.Time
	Flags       []string
	Attachments []Attachment
}

type Attachment struct {
	ID       string
	Name     string
	MimeType string
	Size     int64
}
