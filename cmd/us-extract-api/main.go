package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/mdwhatcott/helps"
	"github.com/mdwhatcott/smarty-cli"
	"github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Lmicroseconds)

	inputs := NewInputs()
	inputs.Flags()
	builder := wireup.NewClientBuilder().WithSecretKeyCredential(inputs.AuthID, inputs.AuthToken).WithDebugHTTPOutput()
	client := builder.BuildUSExtractAPIClient()
	lookup := inputs.AssembleLookup()

	if err := client.SendLookup(lookup); err != nil {
		log.Fatal(err)
	}

	log.Println("Formatted Result:")
	fmt.Println(helps.DumpJSON(lookup.Result))
}

/////////////

type Inputs struct {
	*cli.Inputs

	text             string
	html             string // true, false, or blank
	aggressive       bool
	lineBreaks       bool
	addressesPerLine int

	lookup *extract.Lookup
}

func NewInputs() *Inputs {
	return &Inputs{
		Inputs: cli.NewInputs(),
		lookup: new(extract.Lookup),
	}
}

func (this *Inputs) Flags() {
	flag.StringVar(&this.text, "text", "", "The POST body.")
	flag.StringVar(&this.html, "html", "", "The html field (derived when blank, 'true' or 'false').")
	flag.BoolVar(&this.aggressive, "aggressive", false, "The aggressive bool.")
	flag.BoolVar(&this.lineBreaks, "addr_line_breaks", true, "The addr_line_breaks bool.")
	flag.IntVar(&this.addressesPerLine, "addr_per_line", 0, "T:he add_per_line field.")
	this.ParseFlags()
}

func (this *Inputs) AssembleLookup() *extract.Lookup {
	values, _ := url.ParseQuery(this.RawQuery)
	if this.assembleLookupFromQueryString(values) {
		return this.lookup
	}

	if address, _ := url.Parse(this.RawURL); address != nil {
		if this.assembleLookupFromQueryString(address.Query()) {
			return this.lookup
		}
	}

	this.assembleLookupFromFlags()

	if this.lookup.Text == "" {
		log.Fatal("No data provided.")
	}

	return this.lookup
}
func (this *Inputs) assembleLookupFromFlags() {
	this.lookup.Text = this.text
	this.lookup.AddressesPerLine = this.addressesPerLine
	this.lookup.AddressesWithLineBreaks = this.lineBreaks
	this.lookup.Aggressive = this.aggressive
	this.lookup.HTML = extract.HTMLPayload(this.html)
}
func (this *Inputs) assembleLookupFromQueryString(values url.Values) bool {
	this.lookup.Text = this.text
	this.lookup.AddressesPerLine, _ = strconv.Atoi(values.Get("addr_per_line"))
	this.lookup.AddressesWithLineBreaks, _ = strconv.ParseBool(values.Get("addr_line_breaks"))
	this.lookup.Aggressive, _ = strconv.ParseBool(values.Get("aggressive"))
	this.lookup.HTML = extract.HTMLPayload(values.Get("html"))
	return len(values) > 0
}
