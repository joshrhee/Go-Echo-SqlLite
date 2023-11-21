package handler

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mattn/go-sqlite3"

	"github.com/DATA-DOG/go-sqlmock"
	clocklib "github.com/benbjohnson/clock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_SignUp(t *testing.T) {
	now := time.Date(2021, 1, 2, 3, 4, 5, 6, time.UTC)

	type mocks struct {
		dbRecorder sqlmock.Sqlmock
	}

	cases := []struct {
		name string

		reqBody string

		mockFunc func(mocks) func(ctx context.Context, db *sql.DB, req *http.Request)

		expectedCode int
		expectedJSON string
	}{
		{
			name: "Signup Succeed",
			reqBody: `{
			  "nickname":"dog",
			  "username":"identity",
			  "password":"identity_hard_password"
			}`,
			mockFunc: func(m mocks) func(ctx context.Context, db *sql.DB, req *http.Request) {
				return func(ctx context.Context, db *sql.DB, req *http.Request) {
					m.dbRecorder.ExpectExec("INSERT INTO `user`").WithArgs("dog", "identity", "identity_hard_password", now, now).WillReturnResult(sqlmock.NewResult(123456, 1))
				}
			},

			expectedCode: http.StatusOK,
			expectedJSON: `{"user_id":123456}`,
		},
		{
			name: "Signup Failed_Nickname exceed 12 lengths",
			reqBody: `{
			  "nickname":"longernicknamelongernickname",
			  "username":"identity",
			  "password":"identity_hard_password"
			}`,

			expectedCode: http.StatusBadRequest,
			expectedJSON: `{"reason":"Nickname can't exceed 12 lengths."}`,
		},
		{
			name: "Signup Failed_Username exceed 12 lengths",
			reqBody: `{
			  "nickname":"dog",
			  "username":"identityidentityidentity",
			  "password":"identity_hard_password"
			}`,

			expectedCode: http.StatusBadRequest,
			expectedJSON: `{"reason":"Username can't exceed 12 lengths."}`,
		},
		{
			name: "Signup Failed_Password exceed 50 lengths",
			reqBody: `{
			  "nickname":"멍뭉이",
			  "username":"identity",
			  "password":"identityidentityidentityidentityidentityidentityidentity"
			}`,

			expectedCode: http.StatusBadRequest,
			expectedJSON: `{"reason":"Password can't exceed 50 lengths."}`,
		},
		{
			name: "Signup Failed_Username already exist",
			reqBody: `{
			  "nickname":"멍뭉이",
			  "username":"identity",
			  "password":"identity_hard_password"
			}`,

			mockFunc: func(m mocks) func(ctx context.Context, db *sql.DB, req *http.Request) {
				return func(ctx context.Context, db *sql.DB, req *http.Request) {
					m.dbRecorder.ExpectExec("INSERT INTO `user`").WithArgs("멍뭉이", "identity", "identity_hard_password", now, now).WillReturnError(sqlite3.Error{
						Code:         sqlite3.ErrConstraint,
						ExtendedCode: sqlite3.ErrConstraintUnique,
						SystemErrno:  0,
					})
				}
			},

			expectedCode: http.StatusBadRequest,
			expectedJSON: `{"reason":"Username is already exist"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()

			mockClock := clocklib.NewMock()
			mockClock.Set(now)

			mockDB, mockDBRecorder, err := sqlmock.New()
			assert.NoError(t, err)
			defer func() {
				mockDBRecorder.ExpectClose()
				if err := mockDB.Close(); err != nil {
					t.Error(err)
				}
			}()

			req := httptest.NewRequest(http.MethodPost, "/users/sign-up", strings.NewReader(tc.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tc.mockFunc != nil {
				tc.mockFunc(mocks{
					dbRecorder: mockDBRecorder,
				})(c.Request().Context(), mockDB, c.Request())
			}

			handler := SignUp(mockClock, mockDB)

			assert.NoError(t, handler(c))

			assert.NoError(t, mockDBRecorder.ExpectationsWereMet())

			assert.Equal(t, tc.expectedCode, rec.Code)

			body := rec.Body.String()
			if tc.expectedJSON != "" {
				assert.JSONEq(t, tc.expectedJSON, body)
			}
		})
	}
}
