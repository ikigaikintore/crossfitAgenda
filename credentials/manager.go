package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	folderNameEnv      = "env"
	tokenFileName      = "token.json"
	credentialFileName = "credentials.json"
)

var (
	pathTokenFile      = fmt.Sprintf("%s%s%s", folderNameEnv, string(filepath.Separator), tokenFileName)
	pathCredentialFile = fmt.Sprintf("%s%s%s", folderNameEnv, string(filepath.Separator), credentialFileName)
)

var (
	ErrExpiredToken = errors.New("token expired")
)

type (
	Manager struct {
		cred     []byte
		oauthCfg *oauth2.Config
	}
)

func New() *Manager {
	b, err := os.ReadFile(pathCredentialFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	return &Manager{
		cred: b,
	}
}

func createFolderEnv() {
	if _, err := os.Stat(folderNameEnv); os.IsNotExist(err) {
		if err := os.Mkdir(folderNameEnv, fs.ModeDir); err != nil {
			panic(err)
		}
	}
}

func (m *Manager) GetUrlLogin() string {
	tok, _ := tokenFromFile(pathTokenFile)
	if err := m.valid(tok.Expiry); err != nil && errors.Is(err, ErrExpiredToken) {
		return m.oauthCfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	}
	return ""
}

func (m *Manager) SaveToken(ctx context.Context, authCode string) error {
	tok, err := m.oauthCfg.Exchange(ctx, authCode)
	if err != nil {
		log.Printf("Unable to retrieve token from web: %v", err)
		return err
	}
	saveToken(pathTokenFile, tok)
	return nil
}

func (m *Manager) valid(exp time.Time) error {
	now := time.Now()
	if exp.Before(now) {
		return ErrExpiredToken
	}
	return nil
}

func (m *Manager) SetConfigWithScopes(scopes ...string) error {
	cfg, err := google.ConfigFromJSON(m.cred, scopes...)
	if err != nil {
		return fmt.Errorf("setConfigWithScopes creating from cred file error: %w", err)
	}
	m.oauthCfg = cfg
	return nil
}

func (m *Manager) IsTokenExpired() bool {
	tok, _ := tokenFromFile(pathTokenFile)
	if err := m.valid(tok.Expiry); err != nil && errors.Is(err, ErrExpiredToken) {
		return true
	}
	return false
}

func (m *Manager) GetClient(ctx context.Context) *http.Client {
	if m.oauthCfg == nil {
		panic("Cannot get client, first use SetConfigWithScopes function")
	}

	tok, err := tokenFromFile(pathTokenFile)
	if err != nil {
		tok = m.getTokenFromWeb(ctx)
		saveToken(pathTokenFile, tok)
	}

	if err := m.valid(tok.Expiry); err != nil && errors.Is(err, ErrExpiredToken) {
		tok = m.getTokenFromWeb(ctx)
		saveToken(pathTokenFile, tok)
	}

	return m.oauthCfg.Client(context.Background(), tok)
}

func (m *Manager) getTokenFromWeb(ctx context.Context) *oauth2.Token {
	authURL := m.oauthCfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := m.oauthCfg.Exchange(ctx, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Println(err)
		}
	}()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	createFolderEnv()
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err := json.NewEncoder(f).Encode(token); err != nil {
		panic(err)
	}
}
