package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"notes-server/constants"
	"notes-server/interfaces"
	"notes-server/loggers"
	"notes-server/models"
	"notes-server/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type sessionID struct {
	SID string `json:"sid" validate:"required"`
}

func TokenValidation(db interfaces.ILoginRepository, logger *loggers.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			request := make(map[string]interface{})
			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				logger.Warn(ctx, "error in TokenValidation(), error from utils.GetBodyParams()", err)
				utils.WriteHttpFailure(rw, http.StatusUnauthorized, err)
				return
			}
			if _, ok := request["sid"]; !ok {
				err = errors.New("sid missing")
				logger.Warn(ctx, "error in TokenValidation(), sid missing", err)
				utils.WriteHttpFailure(rw, http.StatusUnauthorized, err)
				return
			}
			if _, ok := request["sid"].(string); !ok {
				err = errors.New("sid invalid")
				logger.Warn(ctx, "error in TokenValidation(), sid invalid", err)
				utils.WriteHttpFailure(rw, http.StatusUnauthorized, err)
				return
			}
			claims := &models.Claims{}
			token, err := jwt.ParseWithClaims(request["sid"].(string), claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(viper.GetString(constants.JwtSecretEnvKey)), nil
			})
			if err != nil {
				logger.Warn(ctx, "error in TokenValidation(), error from jwt.ParseWithClaims()", err)
				utils.WriteHttpFailure(rw, http.StatusUnauthorized, err)
				return
			}
			if !token.Valid {
				err = errors.New("invalid token")
				logger.Warn(ctx, "error in TokenValidation(), invalid token", err)
				utils.WriteHttpFailure(rw, http.StatusUnauthorized, err)
				return
			}
			err = db.ValidateUser(ctx, claims.Email, claims.Name)
			if err != nil {
				logger.Warn(ctx, "error in TokenValidation(), error from db.ValidateUser()", err)
				utils.WriteHttpFailure(rw, http.StatusUnauthorized, err)
				return
			}
			delete(request, "sid")
			req, _ := json.Marshal(request)
			r.Body = io.NopCloser(bytes.NewBuffer(req))
			ctx = context.WithValue(r.Context(), constants.EmailCtxKey, claims.Email)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
