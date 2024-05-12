package messaging

import (
	"context"
	"encoding/json"
	"github.com/tejiriaustin/narx_api/env"
	"log"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type (
	FirebaseMessaging struct {
		ApiKey string
		app    *firebase.App
		conf   *env.Environment
	}
	Message struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
)

var _ Messaging = (*FirebaseMessaging)(nil)

func NewFirebaseMessaging(conf *env.Environment) *FirebaseMessaging {
	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}

	opt := option.WithCredentialsJSON([]byte(serviceAccountKeyFilePath))

	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Firebase load error")
	}

	return &FirebaseMessaging{
		app:  app,
		conf: conf,
	}
}

func (f *FirebaseMessaging) Push(to string, msg string) error {
	ctx := context.Background()

	client, err := f.app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	firebaseMessage, err := UnmarshalFirebaseMessage(ctx, msg)
	if err != nil {
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: firebaseMessage.Title,
			Body:  firebaseMessage.Body,
		},
		Token: f.conf.GetAsString(env.FirebaseRegistrationToken),
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully sent message:", response)
	return nil
}

func UnmarshalFirebaseMessage(ctx context.Context, msg string) (Message, error) {
	var message Message
	err := json.Unmarshal([]byte(msg), &message)
	if err != nil {
		return message, err
	}
	return message, nil
}
