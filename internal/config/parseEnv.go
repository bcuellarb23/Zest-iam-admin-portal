// FIXME: This whole package will be deleted, It won't be used anymore, Know we use godotenv
/* package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func ParseEnvironment() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Printf("Environment Variable file(.env) is not present.")
		// useGlobalEnv = false
	}

	setEnvVariable("OKTA_CLIENT_ID", os.Getenv("OKTA_CLIENT_ID"))
	setEnvVariable("OKTA_CLIENT_SECRET", os.Getenv("OKTA_CLIENT_SECRET"))
	setEnvVariable("OKTA_ISSUER", os.Getenv("OKTA_ISSUER"))

	if os.Getenv("OKTA_CLIENT_ID") == "" {
		log.Printf("Could not resolve a CLIENT_ID environsment variable.")
		os.Exit(1)
	}

	if os.Getenv("OKTA_CLIENT_SECRET") == "" {
		log.Printf("Could not resolve a CLIENT_SECRET environsment variable.")
		os.Exit(1)
	}

	if os.Getenv("OKTA_ISSUER") == "" {
		log.Printf("Could not resolve a ISSUER environsment variable.")
		os.Exit(1)
	}
}

func setEnvVariable(env string, current string) {
	if current != "" {
		return
	}

	file, _ := os.Open(".env")
	defer file.Close()

	lookInFile := bufio.NewScanner(file)
	lookInFile.Split(bufio.ScanLines)

	for lookInFile.Scan() {

		parts := strings.Split(lookInFile.Text(), "=")
		key, value := parts[0], parts[1]
		if key == env {
			os.Setenv(key, value)
		}
	}
} */
