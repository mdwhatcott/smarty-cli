package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"

	"github.com/mdwhatcott/helps"

	"github.com/mdwhatcott/smarty-cli"
)

func main() {
	log.SetFlags(log.Lmicroseconds)

	inputs := NewInputs()
	client := wireup.BuildUSStreetAPIClient(
		wireup.SecretKeyCredential(inputs.AuthID, inputs.AuthToken),
		wireup.DebugHTTPOutput(),
	)
	batch := inputs.PopulateBatch()

	if err := client.SendBatch(batch); err != nil {
		log.Fatal(err)
	}

	var candidates []*street.Candidate
	for _, record := range batch.Records() {
		candidates = append(candidates, record.Results...)
	}
	log.Println("Formatted Result:")
	fmt.Println(helps.DumpJSON(candidates))
}

///////////////////

type Inputs struct {
	*cli.Inputs

	addressee         string
	urbanization      string
	street1           string
	street2           string
	lastLine          string
	secondary         string
	city              string
	state             string
	zipCode           string
	inputID           string
	maxCandidateCount int
	matchStrategy     string

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
	flag.StringVar(&this.addressee, "addressee", "", "The Addresses (US Street API)")
	flag.StringVar(&this.urbanization, "urbanization", "", "The Urbanization (US Street API)")
	flag.StringVar(&this.street1, "street", "", "The Street1 (US Street API)")
	flag.StringVar(&this.street2, "street2", "", "The Street2 (US Street API)")
	flag.StringVar(&this.lastLine, "lastline", "", "The LastLine (US Street API)")
	flag.StringVar(&this.secondary, "secondary", "", "The Secondary (US Street API)")
	flag.StringVar(&this.city, "city", "", "The City (US Street API, US ZIP Code API)")
	flag.StringVar(&this.state, "state", "", "The State (US Street API, US ZIP Code API)")
	flag.StringVar(&this.zipCode, "zipcode", "", "The ZIP Code (US Street API, US ZIP Code API)")
	flag.StringVar(&this.inputID, "input_id", "", "The Input ID (US Street API, US ZIP Code API)")
	flag.IntVar(&this.maxCandidateCount, "candidates", 10, "The max candidate count (US Street API)")
	flag.StringVar(&this.matchStrategy, "match", string(street.MatchStrict), "The Match Strategy (US Street API)")
	this.ParseFlags()
}

func (this *Inputs) PopulateBatch() *street.Batch {
	batch := street.NewBatch()

	var lookups []*street.Lookup
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

func (this *Inputs) assembleLookup() *street.Lookup {
	values, _ := url.ParseQuery(this.RawQuery)
	this.assembleLookupFromQueryString(values)
	if this.lookup.Street != "" {
		return this.lookup
	}

	if address, _ := url.Parse(this.RawURL); address != nil {
		this.assembleLookupFromQueryString(address.Query())
	}
	if this.lookup.Street != "" {
		return this.lookup
	}

	this.assembleLookupFromFlags()

	if this.lookup.Street == "" {
		log.Fatal("No street provided.")
	}

	return this.lookup
}
func (this *Inputs) assembleLookupFromQueryString(values url.Values) {
	this.lookup.Street = values.Get("street")
	this.lookup.Street2 = values.Get("street2")
	this.lookup.City = values.Get("city")
	this.lookup.State = values.Get("state")
	this.lookup.ZIPCode = values.Get("zipcode")
	this.lookup.LastLine = values.Get("lastline")
	this.lookup.Addressee = values.Get("addressee")
	this.lookup.Urbanization = values.Get("urbanization")
	this.lookup.Secondary = values.Get("secondary")
	this.lookup.MatchStrategy = street.MatchStrategy(values.Get("match"))
	this.lookup.MaxCandidates, _ = strconv.Atoi(values.Get("candidates"))
}

func (this *Inputs) assembleLookupFromFlags() {
	this.lookup.Addressee = this.addressee
	this.lookup.Urbanization = this.urbanization
	this.lookup.Street = this.street1
	this.lookup.Street2 = this.street2
	this.lookup.LastLine = this.lastLine
	this.lookup.Secondary = this.secondary
	this.lookup.City = this.city
	this.lookup.State = this.state
	this.lookup.ZIPCode = this.zipCode
	this.lookup.InputID = this.inputID
	this.lookup.MaxCandidates = this.maxCandidateCount
	this.lookup.MatchStrategy = street.MatchStrategy(this.matchStrategy)
}
