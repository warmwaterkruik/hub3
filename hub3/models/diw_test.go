package models_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDIW(t *testing.T) {
	diw, err := getDIW("CUIwQjKBhtnWJo36MrJi")
	if err != nil {
		t.Fatal(err)
	}
	expected := "info@delving.eu"
	if diw.EmailTo != expected {
		t.Errorf("unable to retrieve emailTo: got %v want %v", diw.EmailTo, expected)
	}
}

func TestF(t *testing.T) {
	testWithReturns := `
	Can I order this image.

	Kindest regards,

	Melvin
	`
	feedbackRequest := FeedbackRequest{
		EmailFrom:         "me@delving.eu",
		Feedback:          testWithReturns,
		SourceURI:         "https://example.com/index.html?diw-id=brabantcloud_ton-smits-huis_S07222",
		CollectionSpec:    "ton-smits-huis",
		RecordID:          "S07222",
		RecordTitle:       "No Title",
		RecordCreator:     "Ton Smits",
		RequestID:         "CUIwQjKBhtnWJo36MrJi",
		DryRun:            false,
		RecordImageViewed: "https://i.imgur.com/T6yDl1f.jpg",
	}
	b, err := json.Marshal(&feedbackRequest)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%s", b)

	reader := bytes.NewReader(b)
	r, err := http.NewRequest("POST", "/", reader)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(F)
	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	b, err = json.Marshal(&feedbackRequest)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) != len(b) {
		t.Errorf("Marshalled version should be the same with dry run: got %s want %s", body, b)
	}
}
