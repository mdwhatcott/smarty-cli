package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	cli "bitbucket.org/michael-whatcott/smarty-cli"
	"github.com/mdwhatcott/helps"
	"github.com/smartystreets/smartystreets-go-sdk/us-autocomplete-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Lmicroseconds)

	inputs := NewInputs()
	inputs.Flags()
	builder := wireup.NewClientBuilder().WithSecretKeyCredential(inputs.AuthID, inputs.AuthToken).WithDebugHTTPOutput()
	client := builder.BuildUSAutocompleteAPIClient()
	lookup := inputs.AssembleLookup()

	if err := client.SendLookup(lookup); err != nil {
		log.Fatal(err)
	}

	log.Println("Formatted Result:")
	fmt.Println(helps.DumpJSON(lookup.Results))
}

/////////////

type Inputs struct {
	*cli.Inputs

	prefix             string
	suggestions        int
	geolocatePrecision string
	prefer             string
	preferRatio        float64
	cityFilter         string
	stateFilter        string

	lookup *autocomplete.Lookup
}

func NewInputs() *Inputs {
	return &Inputs{
		Inputs: cli.NewInputs(),
		lookup: new(autocomplete.Lookup),
	}
}

func (this *Inputs) Flags() {
	flag.StringVar(&this.prefix, "prefix", "", "The prefix field.")
	flag.StringVar(&this.geolocatePrecision, "geolocate_precision", "city", "The geolocate_precision field (One of 'city', 'state', or 'none'. A value of 'None' will set the geolocate field to false).")
	flag.StringVar(&this.prefer, "prefer", "", "The prefer field.")
	flag.Float64Var(&this.preferRatio, "prefer_ratio", float64(1.0/3.0), "The prefer_ratio field.")
	flag.StringVar(&this.cityFilter, "city_filter", "", "The city_filter field.")
	flag.StringVar(&this.stateFilter, "state_filter", "", "The state_filter field.")
	flag.IntVar(&this.suggestions, "suggestions", 10, "The suggestions field.")
	this.ParseFlags()
}

func (this *Inputs) AssembleLookup() *autocomplete.Lookup {
	values, _ := url.ParseQuery(this.RawQuery)
	this.assembleLookupFromQueryString(values)
	if this.lookup.Prefix != "" {
		return this.lookup
	}

	if address, _ := url.Parse(this.RawURL); address != nil {
		this.assembleLookupFromQueryString(address.Query())
	}
	if this.lookup.Prefix != "" {
		return this.lookup
	}

	this.assembleLookupFromFlags()

	if this.lookup.Prefix == "" {
		log.Fatal("No prefix provided.")
	}

	return this.lookup
}
func (this *Inputs) assembleLookupFromFlags() {
	this.lookup.Prefix = this.prefix
	this.lookup.CityFilter = strings.Split(this.cityFilter, ",")
	this.lookup.StateFilter = strings.Split(this.stateFilter, ",")
	this.lookup.Preferences = strings.Split(this.prefer, ";")
	this.lookup.PreferRatio = this.preferRatio
	if this.geolocatePrecision == "" {
		this.lookup.Geolocation = autocomplete.GeolocateCity
	} else if this.geolocatePrecision == "state" {
		this.lookup.Geolocation = autocomplete.GeolocateState
	} else if this.geolocatePrecision == "none" {
		this.lookup.Geolocation = autocomplete.GeolocateNone
	}
	this.lookup.MaxSuggestions = this.suggestions
}
func (this *Inputs) assembleLookupFromQueryString(values url.Values) {
	this.lookup.Prefix = values.Get("prefix")
	this.lookup.CityFilter = strings.Split(values.Get("city_filter"), ",")
	this.lookup.StateFilter = strings.Split(values.Get("state_filter"), ",")
	this.lookup.Preferences = strings.Split(values.Get("prefer"), ";")
	this.lookup.PreferRatio, _ = strconv.ParseFloat(values.Get("prefer_ratio"), 64)
	this.lookup.MaxSuggestions, _ = strconv.Atoi(values.Get("suggestions"))
	precision := values.Get("geolocate_precision")
	if precision == "" {
		this.lookup.Geolocation = autocomplete.GeolocateCity
	} else if precision == "state" {
		this.lookup.Geolocation = autocomplete.GeolocateState
	} else if precision == "none" {
		this.lookup.Geolocation = autocomplete.GeolocateNone
	}
}
