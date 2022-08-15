package update_note

type UpdateNoteIn struct {
	NoteId  int    `json:"note_id"`
	NewNote string `json:"new_note"`
}
