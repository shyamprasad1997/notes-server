package controllers

import (
	"net/http"
	"notes-server/models"
	"notes-server/utils"
)

func (c *NotesController) GetNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response, err := c.service.GetNotes(ctx)
	if err != nil {
		c.logger.Warn(ctx, "error in c.service.Login()", err)
		utils.WriteHttpFailure(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteHttpSuccess(w, http.StatusOK, response)
}

func (c *NotesController) AddNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request models.AddNoteRequest
	err := utils.GetBodyParams(r, &request)
	if err != nil {
		c.logger.Warn(ctx, "invalid request", err)
		utils.WriteHttpFailure(w, http.StatusBadRequest, err)
		return
	}
	reponse, err := c.service.AddNote(ctx, request)
	if err != nil {
		c.logger.Warn(ctx, "error in c.service.AddNote()", err)
		utils.WriteHttpFailure(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteHttpSuccess(w, http.StatusCreated, reponse)
}

func (c *NotesController) DeleteNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request models.DeleteNoteRequest
	err := utils.GetBodyParams(r, &request)
	if err != nil {
		c.logger.Warn(ctx, "invalid request", err)
		utils.WriteHttpFailure(w, http.StatusBadRequest, err)
		return
	}
	err = c.service.DeleteNote(ctx, request)
	if err != nil {
		c.logger.Warn(ctx, "error in c.service.DeleteNote()", err)
		utils.WriteHttpFailure(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteHttpSuccess(w, http.StatusOK, "succesfully deleted")
}
