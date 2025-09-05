package workers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Judge0Submission struct {
	LanguageID int    `json:"language_id"`
	SourceCode string `json:"source_code"`
	Stdin      string `json:"stdin,omitempty"`
}

func CreateBatchSubmission(submissionID, sourceCode string, languageID int, testcases []map[string]string) ([]string, error) {
	// Step 1: Build submissions array
	fmt.Println("[Step 1] Building submissions payload...")
	var submissions []Judge0Submission
	for _, tc := range testcases {
		submissions = append(submissions, Judge0Submission{
			LanguageID: languageID,
			SourceCode: sourceCode,
			Stdin:      tc["input"],
		})
	}
	fmt.Printf("[Step 1] Built %d submissions: %+v\n", len(submissions), submissions)

	if len(submissions) == 0 {
		return nil, errors.New("no testcases provided for batch submission")
	}

	// Step 2: Wrap submissions in payload
	payload := map[string]interface{}{
		"submissions": submissions,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal submissions: %v", err)
	}
	fmt.Println("[Step 2] JSON payload ready:", string(data))

	// Step 3: Get URI and API key
	judge0URI := os.Getenv("JUDGE0_URI")
	if judge0URI == "" {
		return nil, errors.New("JUDGE0_URI not set in environment")
	}
	apiKey := os.Getenv("JUDGE0_API_KEY")
	if apiKey == "" {
		return nil, errors.New("JUDGE0_API_KEY not set in environment")
	}
	fmt.Println("[Step 3] Using JUDGE0_URI:", judge0URI)
	fmt.Println("[Step 4] Using JUDGE0_API_KEY:", apiKey)

	// Step 5: Create HTTP request
	fmt.Println("[Step 5] Creating HTTP request...")
	req, err := http.NewRequest("POST", judge0URI, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Step 6: Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RapidAPI-Key", apiKey)
	req.Header.Set("X-RapidAPI-Host", "judge0-ce.p.rapidapi.com")
	fmt.Println("[Step 6] Headers set:", req.Header)

	// Step 7: Send HTTP request
	client := &http.Client{}
	fmt.Println("[Step 7] Sending HTTP request to Judge0...")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()
	fmt.Println("[Step 8] Received response with status:", resp.Status)

	// Step 9: Read raw response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	fmt.Println("[Step 9] Judge0 raw response:", string(bodyBytes))

	// Step 10: Decode response as an array of submissions
	fmt.Println("[Step 10] Decoding JSON response...")
	var respData []struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(bodyBytes, &respData); err != nil {
		fmt.Println("[Step 11] Failed to decode response:", err)
		return nil, fmt.Errorf("failed to decode response JSON: %v", err)
	}

	if len(respData) == 0 {
		fmt.Println("[Step 12] No tokens returned from Judge0")
		return nil, errors.New("no tokens returned from Judge0")
	}

	// Step 13: Extract tokens
	tokens := make([]string, len(respData))
	for i, t := range respData {
		tokens[i] = t.Token
	}
	fmt.Println("[Step 13] Extracted tokens:", tokens)

	return tokens, nil
}
