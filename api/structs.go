package api

// Country endpoint

// Contains the input data from the mmg api, used twice for requests with a date range
// Once for data with cases and one for data with recovered statistics
type InputDataWithScope struct {
	All struct {
		Country    		string                  `json:"country"`
		Population 		int                     `json:"population"`
		Continent 		string                  `json:"continent"`
		Dates      		map[string]interface{}  `json:"dates"`
	}
}

// For the input data of the request without a defined scope in the country endpoint
type InputDataWithoutScope struct {
	All struct {
		Confirmed		int					   	`json:"confirmed"`
		Recovered		int					   	`json:"recovered"`
		Country			string     				`json:"country"`
		Continent		string					`json:"continent"`
		Population		int						`json:"population"`
	}
}

// For the final output in the country endpoint
type OutputData struct {
	Country              string  				`json:"country"`
	Continent            string  				`json:"continent"`
	Scope                string  				`json:"scope"`
	Confirmed            int     				`json:"confirmed"`
	Recovered            int    				`json:"recovered"`
	PopulationPercentage float64 				`json:"population_percentage"`
}

// Stringency endpoint

// For the country name and alpha 3 code retrieved from the restcountries api
type CountryInfo []struct {
	Name       string `json:"name"`
	Alpha3Code string `json:"alpha3Code"`
}

// For the final output in the stringency endpoint
type OutputStrData struct {
	Country 			string  `json:"country"`
	Scope               string  `json:"scope"`
	Stringency			float64 `json:"stringency"`
	Trend				float64 `json:"trend"`
}

// For the input data of the request to the "covidtrackerapi", potentially used for caching data later
/*
type StringencyData struct {
	// Might need to chance countries type
	// Use for checking the if alpha 3 code is reported by the api before doing anything else
	// Potentially just use the Data values
	Countries			[]string 				`json:"countries"`
	Data				map[string]interface{}  `json:"data"`
}
*/

// Stringency data retrieved in the stringency endpoint from the covidtracker api
type StringencyData struct {
	Data struct {
		DateValue        string  `json:"date_value"`
		CountryCode      string  `json:"country_code"`
		Confirmed        int     `json:"confirmed"`
		Deaths           int     `json:"deaths"`
		StringencyActual float64 `json:"stringency_actual"`
		Stringency       float64 `json:"stringency"`
		Message			 string	 `json:"msg"`
	} `json:"stringencyData"`
}

/*
diag represents the diagnostic feedback structure.
It is of the form:
{
	"mmediagroupapi":"<http status code for mmediagroupapi API>",
   	"covidtrackerapi":"<http status code for covidtrackerapi API>",
   	"registered":<number of registered webhooks>,
   	"version":"v1",
   	"uptime":<time in seconds from the last service restart>
}
*/
type Diagnostic struct {
	Mmediagroupapi      	string 		`json:"mmediagroupapi"`
	Covidtrackerapi       	string  	`json:"covidtrackerapi"`
	Register				int			`json:"registered"`
	Version 				string 		`json:"version"`
	Uptime 					string 		`json:"uptime"`
}

// Notifications endpoint

/*
For the webhook data, id is created server side, rest is sent clientside as a post request
It is of the form:
{
	"Id":"<unique identifier>"
	"Url":"<url to invoke on timeout (&& change)>"
	"Timeout":"<time between updates>"
	"Field":"<recovered or confirmed>"
	"Country":"<country name>"
	"Trigger":"<ON_CHANGE or ON_TIMEOUT>" either both change and timeout, or just timeout required
}
*/
type WebhookData struct {
	Id			string	`json:"id"`
	Url			string	`json:"url"`
	Timeout		int		`json:"timeout"`
	Field		string	`json:"field"`
	Country		string	`json:"country"`
	Trigger		string	`json:"trigger"`
}