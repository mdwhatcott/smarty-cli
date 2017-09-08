package cli

import (
	"flag"
	"os"
)

type Inputs struct {
	AuthID    string
	AuthToken string
	RawText   string
	RawQuery  string
	RawURL    string
}

func NewInputs() *Inputs {
	this := &Inputs{}
	this.flags()
	return this
}

func (this *Inputs) flags() {
	flag.StringVar(&this.AuthID, "auth-id", os.Getenv("SMARTY_AUTH_ID"),
		"The auth-id value. Defaults to `SMARTY_AUTH_ID` environment variable value if set.")
	flag.StringVar(&this.AuthToken, "auth-token", os.Getenv("SMARTY_AUTH_TOKEN"),
		"The auth-token value. Defaults to `SMARTY_AUTH_TOKEN` environment variable value if set.")

	flag.StringVar(&this.RawText, "raw", "", "The POST body (US Street API, US ZIP Code API, US Extract API).")
	flag.StringVar(&this.RawQuery, "query", "", "A query string with input values."+authDisclaimerSuffix)
	flag.StringVar(&this.RawURL, "url", "", "A url with query string input values."+authDisclaimerSuffix)
}

const authDisclaimerSuffix = "Even when present, auth-id and auth-token query string values will be ignored. " +
	"(US Street API, US ZIP Code API, US Autocomplete API, US Extract API, International Street API)"
