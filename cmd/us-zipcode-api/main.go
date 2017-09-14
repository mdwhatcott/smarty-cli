package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/mdwhatcott/helps"
	"github.com/mdwhatcott/smarty-cli"
	"github.com/smartystreets/smartystreets-go-sdk/us-zipcode-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Lmicroseconds)

	inputs := NewInputs()
	inputs.Flags()
	builder := wireup.NewClientBuilder().WithSecretKeyCredential(inputs.AuthID, inputs.AuthToken).WithDebugHTTPOutput()
	client := builder.BuildUSZIPCodeAPIClient()
	batch := inputs.PopulateBatch()

	if err := client.SendBatch(batch); err != nil {
		log.Fatal(err)
	}

	var results []*zipcode.Result
	for _, record := range batch.Records() {
		results = append(results, record.Result)
	}
	log.Println("Formatted Result:")
	fmt.Println(helps.DumpJSON(results))
}

/////////////

type Inputs struct {
	*cli.Inputs

	city    string
	state   string
	zipCode string

	lookup *zipcode.Lookup
}

func NewInputs() *Inputs {
	return &Inputs{
		Inputs: cli.NewInputs(),
		lookup: new(zipcode.Lookup),
	}
}

func (this *Inputs) Flags() {
	flag.StringVar(&this.city, "city", "", "The City (US Street API, US ZIP Code API)")
	flag.StringVar(&this.state, "state", "", "The State (US Street API, US ZIP Code API)")
	flag.StringVar(&this.zipCode, "zipcode", "", "The ZIP Code (US Street API, US ZIP Code API)")
	this.ParseFlags()
}

func (this *Inputs) PopulateBatch() *zipcode.Batch {
	batch := zipcode.NewBatch()

	var lookups []*zipcode.Lookup
	json.Unmarshal([]byte(this.RawText), &lookups)
	for _, lookup := range lookups {
		batch.Append(lookup)
	}

	if batch.Length() > 0 {
		return batch
	}

	lookup := this.assembleLookup()
	batch.Append(lookup)
	return batch
}
func (this *Inputs) assembleLookup() *zipcode.Lookup {
	values, _ := url.ParseQuery(this.RawQuery)
	this.assembleLookupFromQueryString(values)
	if this.lookup.City != "" || this.lookup.State != "" || this.lookup.ZIPCode != "" {
		return this.lookup
	}

	if address, _ := url.Parse(this.RawURL); address != nil {
		this.assembleLookupFromQueryString(address.Query())
	}
	if this.lookup.City != "" || this.lookup.State != "" || this.lookup.ZIPCode != "" {
		return this.lookup
	}

	this.assembleLookupFromFlags()

	if this.lookup.City == "" && this.lookup.State == "" && this.lookup.ZIPCode == "" {
		log.Fatal("No data provided.")
	}

	return this.lookup
}
func (this *Inputs) assembleLookupFromFlags() {
	this.lookup.City = this.city
	this.lookup.State = this.state
	this.lookup.ZIPCode = this.zipCode
}
func (this *Inputs) assembleLookupFromQueryString(values url.Values) {
	this.lookup.City = values.Get("city")
	this.lookup.State = values.Get("state")
	this.lookup.ZIPCode = values.Get("zipCode")
}
