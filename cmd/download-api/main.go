package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/mdwhatcott/smarty-cli"
	"github.com/sgreben/flagvar"
)

const (
	USStreetAPI             = "us-street-api"
	USStreetData            = "us-street-data"
	USZIPCodeAPI            = "us-zipcode-api"
	USZIPCodeData           = "us-zipcode-data"
	USAutocompleteAPI       = "us-autocomplete-api"
	USAutocompleteData      = "us-autocomplete-data"
	USExtractAPI            = "us-extract-api"
	InternationalStreetAPI  = "international-street-api"
	InternationalStreetData = "international-street-data"
)

const extension = ".tar.gz"

var targets = map[string]string{
	USStreetAPI:             "us-street-api/linux-amd64",
	USStreetData:            "us-street-api/data",
	USZIPCodeAPI:            "us-zipcode-api/linux-amd64",
	USZIPCodeData:           "us-zipcode-api/data",
	USAutocompleteAPI:       "us-autocomplete-api/linux-amd64",
	USAutocompleteData:      "us-autocomplete-api/data",
	USExtractAPI:            "us-extract-api/linux-amd64",
	InternationalStreetAPI:  "international-street-api/linux-amd64",
	InternationalStreetData: "international-street-api/data",
}

func main() {
	log.SetFlags(log.Lmicroseconds)

	var outputPath string
	var version string
	choices := &flagvar.Enum{Choices: []string{
		USStreetAPI,
		USStreetData,
		USZIPCodeAPI,
		USZIPCodeData,
		USAutocompleteAPI,
		USAutocompleteData,
		USExtractAPI,
		InternationalStreetAPI,
		InternationalStreetData,
	}}
	flag.Var(choices, "package", "Which package? "+choices.Help())
	flag.StringVar(&version, "version", "latest", "Which version?")
	flag.StringVar(&outputPath, "output", "", "Output file path.")
	input := cli.NewInputs()
	input.ParseFlags()

	if choices.String() == "" {
		log.Fatal("Required: -package")
	}
	if outputPath == "" {
		outputPath = choices.String() + extension
	}
	fmt.Println(choices.String(), targets[choices.String()])

	address, err := url.Parse(
		"https://download.api.smartystreets.com" + "/" +
			targets[choices.String()] + "/" +
			version + extension,
	)
	if err != nil {
		log.Fatal(err)
	}
	query := address.Query()
	query.Set("auth-id", input.AuthID)
	query.Set("auth-token", input.AuthToken)
	address.RawQuery = query.Encode()
	request, err := http.NewRequest("GET", address.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	log.Println("Sending download request to:", request.URL)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		log.Fatal("Non-OK status code:", response.Status)
	}
	defer response.Body.Close()

	log.Println("Creating output file...")
	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.Println("Writing output file...")
	n, err := io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Wrote %d bytes to: %s", n, outputPath)
	}
}
