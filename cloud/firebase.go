package cloud

import (
	"OnTrek/utils"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type serviceAccountCreds struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

const firebaseBucketName = "ontrek-99865.firebasestorage.app"

func InitFirebase() (*storage.Client, *utils.StorageConfig, error) {
	ctx := context.Background()
	var credsJSON []byte
	var err error

	credsBase64 := os.Getenv("FIREBASE_CREDENTIALS_BASE64")

	if credsBase64 != "" {
		credsJSON, err = base64.StdEncoding.DecodeString(credsBase64)
		if err != nil {
			return nil, nil, fmt.Errorf("errore decodifica base64: %v", err)
		}
		log.Println("Firebase: Usando credenziali da Variabile d'Ambiente")
	} else {
		filePath := "ontrek-serviceAccountKey.json"
		credsJSON, err = os.ReadFile(filePath)
		if err != nil {
			return nil, nil, fmt.Errorf("nessuna credenziale trovata (nè ENV, nè File): %v", err)
		}
		log.Println("Firebase: Usando credenziali da File Locale")
	}

	var rawCreds serviceAccountCreds
	if err := json.Unmarshal(credsJSON, &rawCreds); err != nil {
		return nil, nil, fmt.Errorf("errore parsing JSON credenziali per config: %v", err)
	}

	storageConfig := &utils.StorageConfig{
		BucketName:  firebaseBucketName,
		ClientEmail: rawCreds.ClientEmail,
		PrivateKey:  []byte(rawCreds.PrivateKey),
	}

	creds, err := google.CredentialsFromJSON(ctx, credsJSON, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, nil, fmt.Errorf("errore creazione credenziali google sdk: %v", err)
	}

	opt := option.WithCredentials(creds)
	firebaseConfig := &firebase.Config{
		StorageBucket: firebaseBucketName,
	}

	app, err := firebase.NewApp(ctx, firebaseConfig, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting storage client: %v", err)
	}

	log.Println("Firebase Storage initialized successfully")

	return client, storageConfig, nil
}

func Middleware(client *storage.Client, config *utils.StorageConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("firebaseStorage", client)

		c.Set("storageConfig", config)

		c.Next()
	}
}
