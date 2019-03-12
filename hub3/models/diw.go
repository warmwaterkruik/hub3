package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"

	html "html/template"
	plaintext "text/template"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type FeedbackRequest struct {
	EmailFrom         string `json:"emailFrom"`
	NameFrom          string `json:"nameFrom"`
	Feedback          string `json:"feedback"`
	SourceURI         string `json:"sourceURI"`
	CollectionSpec    string `json:"collectionSpec"`
	RecordID          string `json:"recordID"`
	RecordTitle       string `json:"recordTitle"`
	RecordCreator     string `json:"recordCreator"`
	RequestID         string `json:"requestID"`
	RecordImageViewed string `json:"recordImageViewed"`
	DryRun            bool   `json:"dryRun"`
}

type DIW struct {
	OrgID    string `firestore:"orgID"`
	EmailTo  string `firestore:"emailTo"`
	Template `firestore:"template"`
}

type Template struct {
	PlainText string `firestore:"plainText"`
	HTML      string `firestore:"html"`
}

func DiwHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Vary", "Origin")
	w.Header().Add("Vary", "Access-Control-Request-Method")
	w.Header().Add("Vary", "Access-Control-Request-Headers")
	// set cors headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "3600")
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	var feedback FeedbackRequest
	err := json.NewDecoder(r.Body).Decode(&feedback)
	if err != nil {
		http.Error(w, "unable to decode body of request", http.StatusBadRequest)
		return
	}
	if !feedback.DryRun {
		diw, err := getDIW(feedback.RequestID)
		if err != nil {
			log.Printf("unable to find requestID %v: %v", feedback.RequestID, err)
			http.Error(w, fmt.Sprintf("unable to find requestID %v", feedback.RequestID), http.StatusBadRequest)
			return
		}

		log.Printf("sending email to %s using SendGrid", diw.EmailTo)
		err = sendEmail(feedback, diw)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to send email: %v", err), http.StatusBadRequest)
			return
		}
	}
	b, err := json.Marshal(feedback)

	if err != nil {
		http.Error(w, "unable to encode body of request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
	return
}

func getDIW(requestID string) (*DIW, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		return nil, err
	}
	diws := client.Collection("diw")
	diwRef, err := diws.Doc(requestID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var diw DIW
	err = diwRef.DataTo(&diw)
	if err != nil {
		return nil, err
	}
	return &diw, nil

}

// sendEmail uses the sendgrid service
func sendEmail(feedback FeedbackRequest, diw *DIW) error {

	plainMessage, err := feedback.renderPlainTextMessage(diw.Template.PlainText)
	if err != nil {
		return err
	}

	htmlMessage, err := feedback.renderHTMLMessage(diw.Template.HTML)
	if err != nil {
		return err
	}

	//
	from := mail.NewEmail("", feedback.EmailFrom)
	subject := fmt.Sprintf("Reactie op: %s (%s)", feedback.RecordTitle, feedback.RecordID)
	to := mail.NewEmail("", diw.EmailTo)
	message := mail.NewSingleEmail(from, subject, to, plainMessage, htmlMessage)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return err
	}
	log.Printf("response: %#v", response)
	log.Printf("send email to %s with id %s", diw.EmailTo, response.Headers["X-Message-Id"])
	return nil
}

func (fr *FeedbackRequest) renderHTMLMessage(tpl string) (string, error) {
	t, err := html.New("").Parse(tpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, fr)
	if err != nil {
		return "", err
	}
	msg := strings.Replace(buf.String(), "\n", "123", -1)
	log.Printf("msg: %#v", msg)
	return msg, nil
}

func (fr *FeedbackRequest) renderPlainTextMessage(tpl string) (string, error) {
	t, err := plaintext.New("").Parse(tpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, fr)
	if err != nil {
		return "", err
	}

	msg := buf.String()
	log.Printf("msg: %#v", msg)
	return msg, nil
}
