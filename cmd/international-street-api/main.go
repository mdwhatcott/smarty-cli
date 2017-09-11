package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	cli "bitbucket.org/michael-whatcott/smarty-cli"
	"github.com/mdwhatcott/helps"
	"github.com/smartystreets/smartystreets-go-sdk/international-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Lmicroseconds)

	inputs := NewInputs()
	builder := wireup.NewClientBuilder().WithSecretKeyCredential(inputs.AuthID, inputs.AuthToken).WithDebugHTTPOutput()
	client := builder.BuildInternationalStreetAPIClient()
	lookup := inputs.AssembleLookup()

	if err := client.SendLookup(lookup); err != nil {
		log.Fatal(err)
	}

	log.Println("Formatted Result:")
	fmt.Println(helps.DumpJSON(lookup.Results))
}

///////////////////

type Inputs struct {
	*cli.Inputs

	country            string
	language           string
	freeform           string
	address1           string
	address2           string
	address3           string
	address4           string
	organization       string
	locality           string
	administrativeArea string
	postalCode         string
	geocode            bool

	lookup *street.Lookup
}

func NewInputs() *Inputs {
	this := &Inputs{
		Inputs: cli.NewInputs(),
		lookup: new(street.Lookup),
	}
	this.flags()
	return this
}

func (this *Inputs) flags() {
	flag.StringVar(&this.country, "country", "", "The country field.")
	flag.StringVar(&this.language, "language", "", "The language field.")
	flag.StringVar(&this.freeform, "freeform", "", "The freeform field.")
	flag.StringVar(&this.address1, "address1", "", "The address1 field.")
	flag.StringVar(&this.address2, "address2", "", "The address2 field.")
	flag.StringVar(&this.address3, "address3", "", "The address3 field.")
	flag.StringVar(&this.address4, "address4", "", "The address4 field.")
	flag.StringVar(&this.organization, "organization", "", "The organization field.")
	flag.StringVar(&this.locality, "locality", "", "The locality field.")
	flag.StringVar(&this.administrativeArea, "administrative_area", "", "The administrative_area field.")
	flag.StringVar(&this.postalCode, "postal_code", "", "The postal_code field.")
	flag.BoolVar(&this.geocode, "geocode", true, "The geocode field.")
	this.ParseFlags()
}

func (this *Inputs) AssembleLookup() *street.Lookup {
	values, _ := url.ParseQuery(this.RawQuery)
	this.assembleLookupFromQueryString(values)
	if this.lookup.Freeform != "" || this.lookup.Address1 != "" {
		return this.lookup
	}

	if address, _ := url.Parse(this.RawURL); address != nil {
		this.assembleLookupFromQueryString(address.Query())
	}
	if this.lookup.Freeform != "" || this.lookup.Address1 != "" {
		return this.lookup
	}

	this.assembleLookupFromFlags()

	if this.lookup.Freeform == "" && this.lookup.Address1 == "" {
		log.Fatal("No address provided.")
	}

	return this.lookup
}
func (this *Inputs) assembleLookupFromQueryString(values url.Values) {
	this.lookup.Country = values.Get("street")
	this.lookup.Language = street.Language(values.Get("language"))
	this.lookup.Organization = values.Get("organization")
	this.lookup.Freeform = values.Get("freeform")
	this.lookup.Address1 = values.Get("address1")
	this.lookup.Address2 = values.Get("address2")
	this.lookup.Address3 = values.Get("address3")
	this.lookup.Address4 = values.Get("address4")
	this.lookup.Locality = values.Get("locality")
	this.lookup.AdministrativeArea = values.Get("administrative_area")
	this.lookup.PostalCode = values.Get("postal_code")
	this.lookup.Geocode = values.Get("geocode") == "true"
}

func (this *Inputs) assembleLookupFromFlags() {
	this.lookup.Country = this.country
	this.lookup.Language = street.Language(this.language)
	this.lookup.Organization = this.organization
	this.lookup.Freeform = this.freeform
	this.lookup.Address1 = this.address1
	this.lookup.Address2 = this.address2
	this.lookup.Address3 = this.address3
	this.lookup.Address4 = this.address4
	this.lookup.Locality = this.locality
	this.lookup.AdministrativeArea = this.administrativeArea
	this.lookup.PostalCode = this.postalCode
	this.lookup.Geocode = this.geocode
}
