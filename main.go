package main

import (
	"crypto/rand"
	"flag"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	"github.com/nstehr/banksy/banksy"
)

//TODO: move tokens to env.var

var (
	port     = flag.Int("port", 8090, "Set the port the service is listening to.")
	genToken = flag.Bool("genToken", false, "Generate the token to secure the hook")
	baseURL  = flag.String("baseUrl", "", "Base url.  Used if working with enterprise github")
	rules    = flag.String("rules", "", "Path to rules yaml file")
)

func main() {
	flag.Parse()
	// provide the helper functionality of generating a token that can be used to secure the webhook
	if *genToken {
		token, err := generateRandomASCIIString(20)
		if err != nil {
			log.Println("Error generating secure token: ", err)
			return
		}
		log.Printf("Secure Token: %s\n", token)
		return
	}
	yamlFile, err := ioutil.ReadFile(*rules)
	if err != nil {
		log.Fatal("Error reading rules yaml file", err)
	}
	l := banksy.NewLabeller(yamlFile)
	apiToken := os.Getenv("GITHUB_API_TOKEN")
	if apiToken == "" {
		log.Fatal("Must specify Github API token using GITHUB_API_TOKEN")
	}
	hookToken := os.Getenv("WEBHOOK_TOKEN")
	if hookToken == "" {
		log.Fatal("Must specify token used to secure webhook using WEBHOOK_TOKEN")
	}
	banksyServer, err := banksy.NewServer(l, *port, hookToken, apiToken, *baseURL)
	if err != nil {
		log.Fatal("Error starting banksy: ", err)
	}
	log.Println("Starting banksy, the github artist")
	banksyServer.Start()
}

// from: https://gist.github.com/denisbrodbeck/635a644089868a51eccd6ae22b2eb800
func generateRandomASCIIString(length int) (string, error) {
	result := ""
	for {
		if len(result) >= length {
			return result, nil
		}
		num, err := rand.Int(rand.Reader, big.NewInt(int64(127)))
		if err != nil {
			return "", err
		}
		n := num.Int64()
		// Make sure that the number/byte/letter is inside
		// the range of printable ASCII characters (excluding space and DEL)
		if n > 32 && n < 127 {
			result += string(n)
		}
	}
}
