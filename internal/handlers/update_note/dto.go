package update_note

type UpdateNoteIn struct {
	NoteId         int    `json:"note_id"`
	NewNoteCaption string `json:"new_note_caption"`
	NewNote        string `json:"new_note"`
}
