package add_note

type AddNoteIn struct {
	NoteCaption string `json:"note_caption"`
	Note        string `json:"note"`
}

type AddNoteOut struct {
	NoteId int `json:"note_id"`
}
