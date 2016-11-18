package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	secretFilePath = "./secret.json"
	codeFile       = "./code_file"
	instaAPI       = "https://api.instagram.com"
	instaAUTH      = "/oauth/authorize/"
	scope          = "basic+public_content+comments"
)

// Secret holds information from json file secret.json
type Secret struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

func readSecretFile(filePath string) (Secret, error) {
	var secret Secret
	secretFile, err := os.Open(filePath)
	if err != nil {
		return Secret{}, err
	}
	defer secretFile.Close()

	if err := json.NewDecoder(secretFile).Decode(&secret); err != nil {
		return Secret{}, err
	}
	return secret, nil
}

func authURL(secret Secret) string {
	format := "%s%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s"
	return fmt.Sprintf(format, instaAPI, instaAUTH, secret.ClientID,
		secret.RedirectURI, scope)
}

func getCodeFromConsole(authurl string) (string, error) {
	fmt.Printf("Please, insert this URL in your browser:\n\t%s\n", authurl)
	fmt.Println("Paste the code")
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	return string(line), nil
}

func readCodeFromFile() (string, bool) {
	f, err := os.Open(codeFile)
	if err != nil {
		return "", false
	}
	line, _, err := bufio.NewReader(f).ReadLine()
	if err != nil {
		log.Printf("[ERROR] instants - %v", err)
		return "", false
	}
	return string(line), true
}

func getCodeFromInsta(s Secret) (string, error) {
	auth := authURL(s)
	code, err := getCodeFromConsole(auth)
	if err != nil {
		return "", err
	}
	return code, nil
}

func getCode(s Secret) (string, error) {
	var err error
	code, ok := readCodeFromFile()
	if !ok {
		if code, err = getCodeFromInsta(s); err != nil {
			return "", err
		}
	}
	if err = saveCode(code); err != nil {
		return "", err
	}
	return code, nil
}

func saveCode(code string) error {
	if code == "" {
		return fmt.Errorf("saveCode: code is empty")
	}

	f, err := os.Create(codeFile)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	if _, err := w.WriteString(code); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func main() {
	s, err := readSecretFile(secretFilePath)
	if err != nil {
		log.Fatalf("[FATAL] instants - %v", err)
	}
	code, err := getCode(s)
	if err != nil {
		log.Fatalf("[FATAL] instants - %v", err)
	}
	fmt.Println(code)
}
