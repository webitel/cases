package tests

import (
	"testing"

	err "github.com/webitel/cases/internal/error"
)

func TestNewInternalError(t *testing.T) {
	id := "internal.error"
	details := "An internal server error occurred"

	err := err.NewInternalError(id, details)

	if err.GetId() != id {
		t.Errorf("expected ID %s, got %s", id, err.GetId())
	}
	if err.GetDetailedError() != details {
		t.Errorf("expected details %s, got %s", details, err.GetDetailedError())
	}
	if err.GetStatusCode() != 500 {
		t.Errorf("expected status code 500, got %d", err.GetStatusCode())
	}
}

func TestNewNotFoundError(t *testing.T) {
	id := "not.found.error"
	details := "The requested resource was not found"

	err := err.NewNotFoundError(id, details)

	if err.GetId() != id {
		t.Errorf("expected ID %s, got %s", id, err.GetId())
	}
	if err.GetDetailedError() != details {
		t.Errorf("expected details %s, got %s", details, err.GetDetailedError())
	}
	if err.GetStatusCode() != 404 {
		t.Errorf("expected status code 404, got %d", err.GetStatusCode())
	}
}

func TestNewBadRequestError(t *testing.T) {
	id := "bad.request.error"
	details := "The request was invalid"

	err := err.NewBadRequestError(id, details)

	if err.GetId() != id {
		t.Errorf("expected ID %s, got %s", id, err.GetId())
	}
	if err.GetDetailedError() != details {
		t.Errorf("expected details %s, got %s", details, err.GetDetailedError())
	}
	if err.GetStatusCode() != 400 {
		t.Errorf("expected status code 400, got %d", err.GetStatusCode())
	}
}

func TestApplicationError_Error(t *testing.T) {
	id := "application.error"
	details := "An application error occurred"
	err := err.NewInternalError(id, details)

	expectedError := "application.error: An application error occurred"
	if err.Error() != expectedError {
		t.Errorf("expected error string '%s', got '%s'", expectedError, err.Error())
	}
}

func TestApplicationError_Translate(t *testing.T) {
	translateFunc := func(id string, params ...interface{}) string {
		if id == "translate.error" {
			return "translated error"
		}
		return id
	}
	err.AppErrorInit(translateFunc)

	id := "translate.error"
	details := "This error needs translation"
	err := err.NewInternalError(id, details)

	err.Translate(translateFunc)

	expectedTranslation := "translated error"
	if err.GetDetailedError() != expectedTranslation {
		t.Errorf("expected translated details '%s', got '%s'", expectedTranslation, err.GetDetailedError())
	}
}

func TestApplicationError_ToJson(t *testing.T) {
	id := "json.error"
	details := "Error for JSON serialization"
	err := err.NewInternalError(id, details)

	expectedJson := `{"id":"json.error","status":"json.error","detail":"Error for JSON serialization","request_id":"","code":500}`
	if err.ToJson() != expectedJson {
		t.Errorf("expected JSON '%s', got '%s'", expectedJson, err.ToJson())
	}
}
