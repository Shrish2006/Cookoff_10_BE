package utils

import (
	"context"
	"fmt"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/google/uuid"
)

type SubmissionInput struct {
	ID         uuid.UUID
	QuestionID string
	LanguageID int
	SourceCode string
	UserID     string
}
// place holder code -- renove when queue is setup
func SaveSubmission(sub SubmissionInput) error {
	if Queries == nil {
		return fmt.Errorf("DB Queries not initialized")
	}

	qid, err := uuid.Parse(sub.QuestionID)
	if err != nil {
		return fmt.Errorf("invalid QuestionID: %v", err)
	}

	userID, err := uuid.Parse(sub.UserID)
	if err != nil {
		return fmt.Errorf("invalid UserID: %v", err)
	}

	status := "pending"

	ctx := context.Background()
	err = Queries.CreateSubmission(ctx, db.CreateSubmissionParams{
		ID:         sub.ID,
		QuestionID: qid,
		LanguageID: int32(sub.LanguageID),
		SourceCode: sub.SourceCode,
		Status:     &status,
		UserID:     uuid.NullUUID{UUID: userID, Valid: true},
		
	})
	return err
}


func GetTestcases(questionID string) ([]map[string]string, error) {
	if Queries == nil {
		return nil, fmt.Errorf("DB Queries not initialized")
	}

	qid, err := uuid.Parse(questionID)
	if err != nil {
		return nil, fmt.Errorf("invalid questionID: %v", err)
	}

	ctx := context.Background()
	rows, err := Queries.GetTestcases(ctx, qid)
	if err != nil {
		return nil, err
	}

	var testcases []map[string]string
	for _, tc := range rows {
		testcases = append(testcases, map[string]string{"input": tc.Input})
	}
	return testcases, nil
}
