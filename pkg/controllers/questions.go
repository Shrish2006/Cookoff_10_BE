package controllers

import (
	"context"
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Question struct {
	ID               uuid.UUID `json:"id"`
	Description      string    `json:"description"`
	Title            string    `json:"title"`
	Qtype            string    `json:"qType"`
	Isbountyactive   bool      `json:"isBountyActive"`
	InputFormat      []string  `json:"inputFormat"`
	Points           int32     `json:"points"`
	Round            int32     `json:"round"`
	Constraints      []string  `json:"constraints"`
	OutputFormat     []string  `json:"outputFormat"`
	SampleTestInput  []string  `json:"sampleTestInput"`
	SampleTestOutput []string  `json:"sampleTestOutput"`
	Explanation      []string  `json:"explanation"`
}

func CreateQuestion(c echo.Context) error {
	var req db.CreateQuestionParams
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	req.ID = uuid.New()
	if err := utils.Queries.CreateQuestion(context.Background(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, req)
}

func GetQuestion(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "UUID GALAT HAI BHAI"})
	}
	q, err := utils.Queries.GetQuestion(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "question not found"})
	}
	return c.JSON(http.StatusOK, q)
}

func GetAllQuestions(c echo.Context) error {
	questions, err := utils.Queries.GetAllQuestions(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, questions)
}

func UpdateQuestion(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "UUID GALAT HAI BHAI"})
	}
	var req db.UpdateQuestionParams
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	req.ID = id
	if err := utils.Queries.UpdateQuestion(context.Background(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func DeleteQuestion(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "UUID GALAT HAI BHAI"})
	}
	if err := utils.Queries.DeleteQuestion(context.Background(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func ActivateBounty(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "UUID GALAT HAI BHAI"})
	}
	if err := utils.Queries.UpdateQuestionBountyActive(context.Background(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func DeactivateBounty(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "UUID GALAT HAI BHAI"})
	}
	if err := utils.Queries.UpdateQuestionBountyInactive(context.Background(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
