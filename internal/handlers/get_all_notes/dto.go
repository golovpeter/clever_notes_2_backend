package get_all_notes

type GetAllNotesOut struct {
	Notes []Note `json:"notes"`
}

type Note struct {
	NoteId  string `json:"note_id" db:"note_id"`
	Caption string `json:"note_caption" db:"note_caption"`
	Text    string `json:"note" db:"note"`
}
