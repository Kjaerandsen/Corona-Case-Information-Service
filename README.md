# Assignment 2: Corona Case Information Service

## Overview
This is a rest api that provides the client functionality to retrieve information about Corona cases occurring in different regions, as well as associated government responses.

The api uses data from the following api's:
* https://mmediagroup.fr/covid-19
* https://covidtracker.bsg.ox.ac.uk/about-api
* https://restcountries.eu/

## Endpoints

```
/corona/v1/country/
/corona/v1/policy/
/corona/v1/diag/
/corona/v1/notifications/
```

The api runs by default on port 8080

The supported request/response pairs are specified in the following.

For the specifications, the following syntax applies: 
* ```{:value}``` indicates mandatory input parameters specified by the user (i.e., the one using *your* service).
* ```{value}``` indicates optional input specified by the user (i.e., the one using *your* service), where `value' can itself contain further optional input. The same notation applies for HTTP parameter specifications (e.g., ```{?param}```).

## Covid-19 Cases per Country

The initial endpoint focuses on return the number of confirmed and recovered Covid-19 cases for a given country. Optionally, the user can specify a date range. 
 * Where such range is specified (in YYYY-MM-DD format), the endpoint provides the newly reported confirmed and recovered cases within this time frame (i.e., not including previous ones). 
 * Where no further constraints/parameters are specified, the endpoint reports the *total numbers* for both confirmed and recovered cases.

### - Request

```
Method: GET
Path: /corona/v1/country/{:country_name}{?scope=begin_date-end_date}
```

```{:country_name}``` refers to the English name for the country as supported by https://mmediagroup.fr/covid-19.

```{?scope=begin_date-end_date}``` indicates the begin date (i.e., the earliest date) of the exchange rate and the end date (i.e., the latest date of the range) of the period for which new cases are reported.

Example request: ```/corona/v1/country/Norway?scope=2020-12-01-2021-01-31```

### - Response

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise.

Body (Example):
```
{
    "country": "Norway",
    "continent": "Europe",
    "scope": "2020-12-01-2021-01-31",
    "confirmed": 13163,
    "recovered": 0,
    "population_percentage": 0.25
}
```

Notes:
* If no date range is provided, the `scope` field reads `"total"`.
* `population_percentage` is the percentage of confirmed cases (in a time frame, if given) within the total population, with two decimal places.

## Covid-19 Policy Stringency Trends

The second endpoint provides an overview of the *current stringency level* of policies regarding Covid-19 for given countries. 
The stringency information is provided by the `stringency_actual` field in the https://covidtracker.bsg.ox.ac.uk API. 
Where the `stringency_actual` field is not filled, the api falls back to the `stringency` field.

### - Request

```
Method: GET
Path: /corona/v1/policy/{:country_name}{?scope=begin_date-end_date}

```{:country_name}``` refers to the English name for the country as supported by https://mmediagroup.fr/covid-19.

```{?scope=begin_date-end_date}``` indicates the first date for which the policy stringency is considered; the end date is the last one. 

Example request: ```/corona/v1/policy/France?scope=2020-12-01-2021-01-31```

### - Response

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise.

Body (Example):
```
{
    "country": "France",
    "scope": "2020-12-01-2021-01-31",
    "stringency": 63.89,
    "trend": -11.11,
}
```

Notes:
* If no date range is provided, the `scope` field should read `"total"`. In this case, the latest `stringency` information should be reported, and `trend` set to 0. 
* If a date range is specified in which either stringency information for begin or end date is not available, report -1 for `stringency`.
* More generally, where information is missing (either stringency information to display, or data for trend calculation), report -1 for `stringency`, and 0 for `trend`.

## Diagnostics interface

The diagnostics interface indicates the availability of all individual services this service depends on. These can be more services than the ones specified above. If you include more, you can specify additional keys with the suffix `api`. The reporting occurs based on status codes returned by the dependent services. The diag interface further provides information about the number of registered webhooks (more details is provided in the next section), and the uptime of the service.

### - Request

```
Method: GET
Path: diag/
```

### - Response

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise. 

Body:
```
{
   "mmediagroupapi": "<http status code for mmediagroupapi API>",
   "covidtrackerapi": "<http status code for covidtrackerapi API>",
   ...
   "registered": <number of registered webhooks>,
   "version": "v1",
   "uptime": <time in seconds from the last service restart>
}
```

Note: ```<some value>``` indicates placeholders for values to be populated by the service.

## Notification Webhook

As an additional feature, users can register webhooks that are triggered by the service based on specified events related to the stringency of policy or confirmed cases. This further includes the specification of country and notification frequency with which the user wishes to be notified.

### Registration of Webhook

### - Request

```
Method: POST
Path: /corona/v1/notifications/
```

The body contains 
 * the URL to be triggered upon event (the service that should be invoked)
 * the frequency with which the invocation occurs in seconds (timeout)
 * the information of interest (`stringency` of policy or `confirmed` cases)
 * indication whether notifications should only be sent when information has changed ("ON_CHANGE") - for example, the stringency has changed since the last call *AND* the timeout is reached -, or be sent in any case whenever the specified timeout expires ("ON_TIMEOUT")
 * the country for which the trigger applies

Body (Example):
```
{
   "url": "https://localhost:8080/client/",
   "timeout": 3600,
   "field": "stringency",
   "country": "France",
   "trigger": "ON_CHANGE"
}
```

**UPDATE: The trigger parameters are `ON_CHANGE` and `ON_TIMEOUT`, not `ON_UPDATE`. This has been fixed now.**

### - Response

The response contains the ID for the registration that can be used to see detail information or to delete the webhook registration. The format of the ID is not prescribed, as long it is unique. Consider best practices for determining IDs.

* Content type: text/plain
* Status code: Choose an appropriate status code

Body (Example):
```
OIdksUDwveiwe
```

### Deletion of Webhook

### - Request

```
Method: DELETE
Path: /corona/v1/notifications/{id}
```

* {id] is the ID for the webhook registration

### - Response

Implement the response according to best practices.

### View registered webhook

### - Request

```
Method: GET
Path: /corona/v1/notifications/{id}
```

* {id] is the ID for the webhook registration

### - Response

The response is similar to the POST request body, but further includes the ID assigned by the server upon adding the webhook.

Body (Example):
```
{
   "id": "OIdksUDwveiwe",
   "url": "http://localhost:8080/client/",
   "timeout": 3600,
   "field": "stringency",
   "country": "France",
   "trigger": "ON_CHANGE"
}
```

### View all registered webhooks

### - Request

```
Method: GET
Path: /corona/v1/notifications/
```

### - Response

The response is a collection of registered webhooks as specified in the POST body, alongside the server-defined ID.

Body (Example):
```
[{
   "id": "OIdksUDwveiwe",
   "url": "https://localhost:8080/client/",
   "timeout": 3600,
   "information": "stringency",
   "country": "France",
   "trigger": "ON_CHANGE"
},
...
]
```

### Webhook Invocation

```
Method: POST
Path: <url specified in the corresponding webhook registration>
```

Any invocation of the registered webhook has the format of the output corresponding to the `information` that has been registered. For example, if `stringency` was specified as information during webhook registration, the structure of the body would follow the policy endpoint output; conversely, if registering `confirmed` as information value for the webhook, the latest (no date ranges) information about confirmed Covid-19 cases for the given country are returned. 

### Credits:

The readme is based on the assignment info retrieved from 
https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/wikis/Assignment-2

Uptime code is based on code from 
https://stackoverflow.com/questions/37992660/golang-retrieve-application-uptime

The url request is based on code from RESTclient found at
"https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/blob/master/RESTclient/cmd/main.go"