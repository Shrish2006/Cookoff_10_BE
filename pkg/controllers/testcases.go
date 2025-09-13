package controllers

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/validator"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateTestCaseRequest struct {
	ExpectedOutput string      `json:"expected_output" validate:"required"`
	Memory         interface{} `json:"memory" validate:"required"`
	Input          string      `json:"input" validate:"required"`
	Hidden         bool        `json:"hidden"`
	Runtime        interface{} `json:"runtime" validate:"required"`
	QuestionID     string      `json:"question_id" validate:"required,uuid"`
}

type UpdateTestCaseRequest struct {
	ExpectedOutput string      `json:"expected_output"`
	Memory         interface{} `json:"memory"`
	Input          string      `json:"input"`
	Hidden         *bool       `json:"hidden"`
	Runtime        interface{} `json:"runtime"`
	QuestionID     string      `json:"question_id" validate:"omitempty,uuid"`
}

func CreateTestCase(c echo.Context) error {
	var req CreateTestCaseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid request body",
			"error":  err.Error(),
		})
	}

	if err := validator.ValidatePayload(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Validation failed",
			"error":  err.Error(),
		})
	}

	questionID, _ := uuid.Parse(req.QuestionID)

	memory, err := utils.InterfaceToNumeric(req.Memory)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid memory value",
			"error":  err.Error(),
		})
	}

	runtime, err := utils.InterfaceToNumeric(req.Runtime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid runtime value",
			"error":  err.Error(),
		})
	}

	testCase, err := utils.Queries.CreateTestCase(c.Request().Context(), db.CreateTestCaseParams{
		ID:             uuid.New(),
		ExpectedOutput: req.ExpectedOutput,
		Memory:         memory,
		Input:          req.Input,
		Hidden:         req.Hidden,
		Runtime:        runtime,
		QuestionID:     questionID,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Failed to create test case",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"status":    "success",
		"test_case": testCase,
	})
}

func GetTestCase(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid test case ID",
			"error":  err.Error(),
		})
	}

	testCase, err := utils.Queries.GetTestCase(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status": "Test case not found",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":    "success",
		"test_case": testCase,
	})
}

func UpdateTestCase(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid test case ID",
			"error":  err.Error(),
		})
	}

	existing, err := utils.Queries.GetTestCase(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status": "Test case not found",
			"error":  err.Error(),
		})
	}

	var req UpdateTestCaseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid request body",
			"error":  err.Error(),
		})
	}

	if req.ExpectedOutput == "" {
		req.ExpectedOutput = existing.ExpectedOutput
	}
	if req.Input == "" {
		req.Input = existing.Input
	}
	if req.Hidden == nil {
		req.Hidden = &existing.Hidden
	}

	questionID := existing.QuestionID
	if req.QuestionID != "" {
		parsedID, err := uuid.Parse(req.QuestionID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status": "Invalid question ID",
				"error":  err.Error(),
			})
		}
		questionID = parsedID
	}

	memory := existing.Memory
	if req.Memory != nil {
		memory, err = utils.InterfaceToNumeric(req.Memory)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status": "Invalid memory value",
				"error":  err.Error(),
			})
		}
	}

	runtime := existing.Runtime
	if req.Runtime != nil {
		runtime, err = utils.InterfaceToNumeric(req.Runtime)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status": "Invalid runtime value",
				"error":  err.Error(),
			})
		}
	}

	updated, err := utils.Queries.UpdateTestCase(c.Request().Context(), db.UpdateTestCaseParams{
		ID:             id,
		ExpectedOutput: req.ExpectedOutput,
		Memory:         memory,
		Input:          req.Input,
		Hidden:         *req.Hidden,
		Runtime:        runtime,
		QuestionID:     questionID,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Failed to update test case",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":    "success",
		"test_case": updated,
	})
}

func DeleteTestCase(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid test case ID",
			"error":  err.Error(),
		})
	}

	_, err = utils.Queries.GetTestCase(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status": "Test case not found",
			"error":  err.Error(),
		})
	}

	err = utils.Queries.DeleteTestCase(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Failed to delete test case",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "Test case deleted successfully",
	})
}

func GetTestCasesByQuestion(c echo.Context) error {
	questionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid question ID",
			"error":  err.Error(),
		})
	}

	testCases, err := utils.Queries.GetTestCasesByQuestion(c.Request().Context(), questionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Failed to fetch test cases",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":      "success",
		"test_cases":  testCases,
		"total_count": len(testCases),
	})
}

func GetPublicTestCasesByQuestion(c echo.Context) error {
	questionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid question ID",
			"error":  err.Error(),
		})
	}

	testCases, err := utils.Queries.GetPublicTestCasesByQuestion(c.Request().Context(), questionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Failed to fetch test cases",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":      "success",
		"test_cases":  testCases,
		"total_count": len(testCases),
	})
}

func GetAllTestCases(c echo.Context) error {
	testCases, err := utils.Queries.GetAllTestCases(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Failed to fetch test cases",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":      "success",
		"test_cases":  testCases,
		"total_count": len(testCases),
	})
}
