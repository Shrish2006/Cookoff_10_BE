package controllers


import (
	"net/http"
	"math/big"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/validator"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/auth"
	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignupRequest struct {
	Email string `json:"email"    validate:"required,email"`
	Name  string `json:"name"     validate:"required"`
	RegNo string `json:"reg_no"   validate:"required"`
	// Key   string `json:"fuck_you" validate:"required"`
}

func Signup(c echo.Context) error {
	var payload SignupRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := validator.ValidatePayload(payload); err != nil {
		return c.JSON(http.StatusNotAcceptable, echo.Map{"error": "invalid input"})
	}

	_, err := utils.Queries.GetUserByEmail(c.Request().Context(), payload.Email)
	if err == nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": "User already exists"})
	}
	
	password := auth.PasswordGenerator(6)
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error" : err.Error()})
	}

	id, err := uuid.NewV7()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	_, err = utils.Queries.CreateUser(c.Request().Context(), db.CreateUserParams{
		ID:             id,
		Email:          payload.Email,
		RegNo:          payload.RegNo,
		Password:       string(hashed),
		Role:           "user",
		RoundQualified: 0,
		Score: pgtype.Numeric{
			Int:              big.NewInt(0),
			Exp:              0,
			NaN:              false,
			InfinityModifier: 0,
			Valid:            true,
		},
		Name:           payload.Name,
	})
	if err != nil {
		return c.JSON(http.StatusCreated, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user added",
		"email": payload.Email,
		"password": password,
	})
}
