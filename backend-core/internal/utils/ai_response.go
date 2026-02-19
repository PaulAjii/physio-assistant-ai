package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func ForwardAudioToAI(filePath, aiBackendUri string) ([]byte, error) {
	// TODO: Fetch audio from path
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %w", err)
	}

	_, err = io.Copy(part, f)
	if err != nil {
		return nil, fmt.Errorf("error copying file to form: %w", err)
	}
	writer.Close()

	// TODO: Send to the Nestjs Backend
	req, err := http.NewRequest(http.MethodPost, aiBackendUri, &body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI backend error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
