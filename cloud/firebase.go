package cloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func InitFirebase() (*storage.Client, error) {
	ctx := context.Background()
	var credsJSON []byte
	var err error

	credsBase64 := os.Getenv("FIREBASE_CREDENTIALS_BASE64")

	if credsBase64 != "" {
		credsJSON, err = base64.StdEncoding.DecodeString(credsBase64)
		if err != nil {
			return nil, fmt.Errorf("errore decodifica base64: %v", err)
		}
		log.Println("Firebase: Usando credenziali da Variabile d'Ambiente")

	} else {
		filePath := "ontrek-serviceAccountKey.json"
		credsJSON, err = os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("nessuna credenziale trovata (nè ENV, nè File): %v", err)
		}
		log.Println("Firebase: Usando credenziali da File Locale")
	}

	creds, err := google.CredentialsFromJSON(ctx, credsJSON, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, fmt.Errorf("errore parsing credenziali google: %v", err)
	}

	opt := option.WithCredentials(creds)

	config := &firebase.Config{
		StorageBucket: "NOME-TUO-BUCKET.appspot.com",
	}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting storage client: %v", err)
	}

	log.Println("Firebase Storage initialized successfully")
	return client, nil
}

func Middleware(client *storage.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("firebaseStorage", client)
		c.Next()
	}
}
