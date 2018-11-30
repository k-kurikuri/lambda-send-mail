package main

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"log"
	"os"
	"time"
)

const (
	rfc2822 = "Mon Jan 02 15:04:05 -0700 2006"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	lambda.Start(Handler)
}

// lambda handler
func Handler() (Response, error) {
	config := getOAuthConfig()
	token := getOAuthToken()
	client := config.Client(context.Background(), &token)

	srv, err := gmail.New(client)
	if err != nil {
		log.Print("Unable to retrieve gmail Client")
		return Response{}, err
	}

	emailContent := getEmailContent()

	message := gmail.Message{
		Raw:      base64.URLEncoding.EncodeToString([]byte(emailContent)),
		LabelIds: []string{"INBOX"},
	}

	m, err := srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Print("Unable to insert messages")
		return Response{}, err
	}

	return Response{
		Message: m.Id + " Gmail Send executed successfully!",
	}, nil
}

// OAuth2に必要な設定情報を返す
func getOAuthConfig() oauth2.Config {
	return oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{os.Getenv("SCOPES")},
	}
}

// OAuth2に必要なTokenを返す
func getOAuthToken() oauth2.Token {
	expiry, _ := time.Parse("2006-01-02", "2018-04-16")

	return oauth2.Token{
		AccessToken:  os.Getenv("ACCESS_TOKEN"),
		TokenType:    os.Getenv("TOKEN_TYPE"),
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		Expiry:       expiry,
	}
}

// Email送信に必要な情報を返す
func getEmailContent() string {
	emailDate := time.Now().Format(rfc2822)

	return "Date: " + emailDate + "\r\n" +
		"From: " + os.Getenv("EMAIL_FROM") + "\r\n" +
		"To: " + os.Getenv("EMAIL_TO") + "\r\n" +
		"Subject: " + os.Getenv("EMAIL_SUBJECT") + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + os.Getenv("EMAIL_BODY")
}
