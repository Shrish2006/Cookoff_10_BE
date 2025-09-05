package controllers

import (
	"fmt"
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/workers"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SubmissionRequest struct {
	SourceCode string `json:"source_code" validate:"required"`
	LanguageID int    `json:"language_id" validate:"required"`
	QuestionID string `json:"question_id" validate:"required"`
	UserID     string `json:"user_id"` // optional
}

func SubmitCode(c echo.Context) error {
	var req SubmissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Generate a new submission ID
	submissionID := uuid.New()

	// Fetch testcases using sqlc Queries
	testcases, err := utils.GetTestcases(req.QuestionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch testcases"})
	}

	// Create batch submission in Judge0
	tokens, err := workers.CreateBatchSubmission(submissionID.String(), req.SourceCode, req.LanguageID, testcases)
	if err != nil {
		fmt.Println("CreateBatchSubmission error:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create batch submission"})
	}

	// Save submission in DB with status "pending"
	sub := utils.SubmissionInput{
		ID:         submissionID,
		QuestionID: req.QuestionID,
		LanguageID: req.LanguageID,
		SourceCode: req.SourceCode,
		UserID:     req.UserID,
	}
	if err := utils.SaveSubmission(sub); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save submission record"})
	}

	// Return submission ID and Judge0 tokens
	return c.JSON(http.StatusOK, map[string]interface{}{
		"submission_id": submissionID,
		"tokens":        tokens,
	})
}
