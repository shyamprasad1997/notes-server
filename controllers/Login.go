package controllers

import (
	"net/http"
	"notes-server/models"
	"notes-server/utils"
)

func (c *LoginController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request models.LoginRequest
	err := utils.GetBodyParams(r, &request)
	if err != nil {
		c.logger.Warn(ctx, "invalid request", err)
		utils.WriteHttpFailure(w, http.StatusBadRequest, err)
		return
	}
	response, err := c.service.Login(ctx, request)
	if err != nil {
		c.logger.Warn(ctx, "error in c.service.Login()", err)
		utils.WriteHttpFailure(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteHttpSuccess(w, http.StatusOK, response)
}

func (c *LoginController) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request models.SignUpRequest
	err := utils.GetBodyParams(r, &request)
	if err != nil {
		c.logger.Warn(ctx, "invalid request", err)
		utils.WriteHttpFailure(w, http.StatusBadRequest, err)
		return
	}
	err = c.service.SignUp(ctx, request)
	if err != nil {
		c.logger.Warn(ctx, "error in c.service.SignUp()", err)
		utils.WriteHttpFailure(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteHttpSuccess(w, http.StatusCreated, "user created")
}
