package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
	autocomplete "github.com/smartystreets/smartystreets-go-sdk/us-autocomplete-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"

	cli "github.com/mdwhatcott/smarty-cli"
)

func main() {
	inputs := NewInputs()
	inputs.Flags()
	client := wireup.BuildUSAutocompleteAPIClient(
		wireup.CustomBaseURL(inputs.baseURL),
		wireup.SecretKeyCredential(inputs.AuthID, inputs.AuthToken),
	)
	lookup := inputs.AssembleLookup()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	line := new(bytes.Buffer)

keyPressListenerLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break keyPressListenerLoop
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if line.Len() == 0 {
					continue
				}
				content := line.String()
				line.Reset()
				line.WriteString(content[:len(content)-1])
			case termbox.KeyEnter:
				line.Reset()
			case termbox.KeySpace:
				line.WriteRune(' ')
			default:
				line.WriteRune(ev.Ch)
			}
		case termbox.EventError:
			panic(ev.Err)
		}

		lookup.Prefix = line.String()
		err := client.SendLookup(lookup)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println()
		fmt.Println(">>>", lookup.Prefix)
		fmt.Println()
		for r, result := range lookup.Results {
			fmt.Printf("%d: [%s]\n", r, result.Text)
		}
		lookup.Results = nil
	}
}

type Inputs struct {
	*cli.Inputs

	baseURL string

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
	flag.StringVar(&this.baseURL, "baseURL", "https://us-autocomplete.api.smartystreets.com", "The URL")
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
