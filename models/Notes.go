package models

type Note struct {
	Id        int32  `json:"id"`
	Note      string `json:"note"`
	CreatedBy string `json:"-"`
}

type AddNoteRequest struct {
	Email string
	Note  string `json:"note" validate:"required"`
}

type AddNoteResponse struct {
	Id int32 `json:"id"`
}

type DeleteNoteRequest struct {
	Id int32 `json:"id" validate:"required"`
}
