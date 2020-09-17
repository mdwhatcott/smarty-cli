package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	reverse "github.com/smartystreets/smartystreets-go-sdk/us-reverse-geo-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"

	"github.com/mdwhatcott/smarty-cli"
	"github.com/mdwhatcott/smarty-cli/helps"
)

func main() {
	log.SetFlags(log.Lmicroseconds)

	inputs := NewInputs()
	client := wireup.BuildUSReverseGeocodingAPIClient(
		wireup.CustomBaseURL(inputs.baseURL),
		wireup.WithLicenses(inputs.Licenses()...),
		wireup.SecretKeyCredential(inputs.AuthID, inputs.AuthToken),
		wireup.DebugHTTPOutput(),
	)
	lookup := inputs.PopulateLookup()

	if err := client.SendLookup(lookup); err != nil {
		log.Fatal(err)
	}

	log.Println("Formatted Result:")
	fmt.Println(helps.DumpJSON(lookup.Response.Results))
}

///////////////////

type Inputs struct {
	*cli.Inputs

	baseURL  string
	licenses string

	latitude  float64
	longitude float64

	lookup *reverse.Lookup
}

func NewInputs() *Inputs {
	this := &Inputs{
		Inputs: cli.NewInputs(),
		lookup: new(reverse.Lookup),
	}
	this.flags()
	return this
}

func (this *Inputs) flags() {
	flag.StringVar(&this.baseURL, "baseURL", os.Getenv("SMARTY_US_REVERSE_GEO_API"), "The URL")
	flag.StringVar(&this.licenses, "licenses", "us-reverse-geocoding-cloud", "The licenses")
	flag.Float64Var(&this.latitude, "latitude", 40.25, "The latitude")
	flag.Float64Var(&this.longitude, "longitude", -111.67, "The longitude")
	this.ParseFlags()
}

func (this *Inputs) Licenses() []string {
	return strings.Split(this.licenses, ",")
}

func (this *Inputs) PopulateLookup() *reverse.Lookup {
	values, _ := url.ParseQuery(this.RawQuery)
	this.assembleLookupFromQueryString(values)

	if this.lookup.Latitude != 0 && this.lookup.Longitude != 0 {
		return this.lookup
	}

	if address, _ := url.Parse(this.RawURL); address != nil {
		this.assembleLookupFromQueryString(address.Query())
	}
	if this.lookup.Latitude != 0 && this.lookup.Longitude != 0 {
		return this.lookup
	}

	this.assembleLookupFromFlags()

	if this.lookup.Latitude == 0 && this.lookup.Longitude == 0 {
		log.Fatal("No street provided.")
	}

	return this.lookup
}
func (this *Inputs) assembleLookupFromQueryString(values url.Values) {
	this.lookup.Latitude = ParseFloat64(values.Get("latitude"))
	this.lookup.Longitude = ParseFloat64(values.Get("longitude"))
}

func (this *Inputs) assembleLookupFromFlags() {
	this.lookup.Latitude = this.latitude
	this.lookup.Longitude = this.longitude
}

func ParseFloat64(raw string) float64 {
	parsed, _ := strconv.ParseFloat(raw, 64)
	return parsed
}
