package handler

import (
	"database/sql"
	"identity-coding-test/internal/database"
	"net/http"

	"github.com/benbjohnson/clock"
	"github.com/labstack/echo/v4"
)

type signUpRequestBody struct {
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type signUpErrorResponse struct {
	Reason string `json:"reason"`
}

type signUpOKResponse struct {
	UserID int64 `json:"user_id"`
}

func SignUp(clock clock.Clock, db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var body signUpRequestBody
		err := c.Bind(&body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &signUpErrorResponse{
				Reason: "bad request",
			})
		}

		if body.Nickname == "" || body.Username == "" || body.Password == "" {
			return c.JSON(http.StatusBadRequest, &signUpErrorResponse{
				Reason: "bad request",
			})

		}

		userID, err := database.CreateUser(c.Request().Context(), db, body.Nickname, body.Username, body.Password, clock.Now())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &signUpErrorResponse{
				Reason: "internal server error",
			})
		}

		return c.JSON(http.StatusOK, &signUpOKResponse{
			UserID: userID,
		})
	}
}
