# Moesif Middleware for Go

[![Built For][ico-built-for]][link-built-for]
[![Software License][ico-license]][link-license]
[![Source Code][ico-source]][link-source]

Go Middleware that logs API Calls and sends to [Moesif](https://www.moesif.com) for API analytics and log analysis.

[Source Code on GitHub](https://github.com/moesif/moesifmiddleware-go)

## How to install
Run the following commands:

`moesifmiddleware-go` can be installed like any other Go library through go get:

```bash
go get github.com/moesif/moesifmiddleware-go
```

Or, if you are already using Go Modules, specify a version number as well:

```bash
go get github.com/moesif/moesifmiddleware-go@v1.0.7
```

## How to use

Add middleware to your application.

```go
import(
    moesifmiddleware "github.com/moesif/moesifmiddleware-go"
)

func handle(w http.ResponseWriter, r *http.Request) {
	// Your API Logic
}

var moesifOptions = map[string]interface{} {
        "Application_Id": "Your Moesif Application Id",
        "Log_Body": true,
}
http.Handle("/api", moesifmiddleware.MoesifMiddleware(http.HandlerFunc(handle), moesifOption))
```

Your Moesif Application Id can be found in the [_Moesif Portal_](https://www.moesif.com/).
After signing up for a Moesif account, your Moesif Application Id will be displayed during the onboarding steps. 

## Optional: Capturing outgoing API calls
In addition to your own APIs, you can also start capturing calls out to third party services via the following method:

```go
moesifmiddleware.StartCaptureOutgoing(moesifOption)
```

#### `handler func(ResponseWriter, *Request)`
(__required__), HandlerFunc registers the handler function for the given pattern.

#### `moesifOption`
(__required__), _map[string]interface{}_, are the configuration options for your application. Please find the details below on how to configure options.

## Configuration options

### __`Application_Id`__
(__required__), _string_, is obtained via your Moesif Account, this is required.
Your Moesif Application Id can be found in the [_Moesif Portal_](https://www.moesif.com/).
After signing up for a Moesif account, your Moesif Application Id will be displayed during the onboarding steps. 

You can always find your Moesif Application Id at any time by logging 
into the [_Moesif Portal_](https://www.moesif.com/), click on the top right menu,
and then clicking _Installation_.

### __`Should_Skip`__
(optional) _(request, response) => boolean_, a function that takes a request and a response,
and returns true if you want to skip this particular event.

### __`Identify_User`__
(optional, but highly recommended) _(request, response) => string_, a function that takes a request and response, and returns a string that is the user id used by your system. While Moesif tries to identify users automatically, but different frameworks and your implementation might be very different, it would be helpful and much more accurate to provide this function.

### __`Identify_Company`__
(optional) _(request, response) => string_, a function that takes a request and response, and returns a string that is the company id for this event.

### __`Get_Metadata`__
(optional) _(request, response) => dictionary_, a function that takes a request and response, and
returns a dictionary (must be able to be encoded into JSON). This allows you
to associate this event with custom metadata. For example, you may want to save a VM instance_id, a trace_id, or a tenant_id with the request.

### __`Get_Session_Token`__
(optional) _(request, response) => string_, a function that takes a request and response, and returns a string that is the session token for this event. Moesif tries to get the session token automatically, but if this doesn't work for your service, you should use this to identify sessions.

### __`Request_Header_Masks`__
(optional) _() => []string_, a function that returns array of strings to mask specific request header fields.

### __`Request_Body_Masks`__
(optional) _() => []string_, a function that returns array of strings to mask specific request body fields.

### __`Response_Header_Masks`__
(optional) _() => []string_, a function that returns array of strings to mask specific response header fields.

### __`Response_Body_Masks`__
(optional) _() => []string_, a function that returns array of strings to mask specific response body fields.

### __`Debug`__
(optional) _boolean_, a flag to see debugging messages.

### __`Log_Body`__
(optional) _boolean_, Default true. Set to false to remove logging request and response body to Moesif.

## Options for logging outgoing calls

The options below are applied to outgoing API calls. The request and response objects passed in are [Request](https://golang.org/src/net/http/request.go) request and [Response](https://golang.org/src/net/http/response.go) response objects.

### __`Should_Skip_Outgoing`__
(optional) _(request, response) => boolean_, a function that takes a request and response, and returns true if you want to skip this particular event.

### __`Identify_User_Outgoing`__
(optional, but highly recommended) _(request, response) => string_, a function that takes request and response, and returns a string that is the user id used by your system. While Moesif tries to identify users automatically,
but different frameworks and your implementation might be very different, it would be helpful and much more accurate to provide this function.

### __`Identify_Company_Outgoing`__
(optional) _(request, response) => string_, a function that takes request and response, and returns a string that is the company id for this event.

### __`Get_Metadata_Outgoing`__
(optional) _(request, response) => dictionary_, a function that takes request and response, and
returns a dictionary (must be able to be encoded into JSON). This allows
to associate this event with custom metadata. For example, you may want to save a VM instance_id, a trace_id, or a tenant_id with the request.

### __`Get_Session_Token_Outgoing`__
(optional) _(request, response) => string_, a function that takes request and response, and returns a string that is the session token for this event. Again, Moesif tries to get the session token automatically, but if you setup is very different from standard, this function will be very help for tying events together, and help you replay the events.

### __`Log_Body_Outgoing`__
(optional) _boolean_, Default true. Set to false to remove logging request and response body to Moesif.

## Update User

### UpdateUser method
Create or update a user profile in Moesif.
The metadata field can be any customer demographic or other info you want to store.
Only the `UserId` field is required.
This method is a convenient helper that calls the Moesif API lib.
For details, visit the [Go API Reference](https://www.moesif.com/docs/api?go#update-a-user).

```go
import (
	moesifmiddleware "github.com/moesif/moesifmiddleware-go"
)

func literalFieldValue(value string) *string {
    return &value
}

var moesifOptions = map[string]interface{} {
	"Application_Id": "Your Moesif Application Id",
	"Log_Body": true,
}

// Campaign object is optional, but useful if you want to track ROI of acquisition channels
// See https://www.moesif.com/docs/api#users for campaign schema
campaign := models.CampaignModel {
  UtmSource: literalFieldValue("google"),
  UtmMedium: literalFieldValue("cpc"), 
  UtmCampaign: literalFieldValue("adwords"),
  UtmTerm: literalFieldValue("api+tooling"),
  UtmContent: literalFieldValue("landing"),
}
  
// metadata can be any custom dictionary
metadata := map[string]interface{}{
  "email": "john@acmeinc.com",
  "first_name": "John",
  "last_name": "Doe",
  "title": "Software Engineer",
  "sales_info": map[string]interface{}{
      "stage": "Customer",
      "lifetime_value": 24000,
      "account_owner": "mary@contoso.com",
  },
}

// Only UserId is required
user := models.UserModel{
  UserId:  "12345",
  CompanyId:  literalFieldValue("67890"), // If set, associate user with a company object
  Campaign:  &campaign,
  Metadata:  &metadata,
}

// Update User
moesifmiddleware.UpdateUser(&user, moesifOption)
```

### UpdateUsersBatch method
Similar to UpdateUser, but used to update a list of users in one batch. 
Only the `UserId` field is required.
This method is a convenient helper that calls the Moesif API lib.
For details, visit the [Go API Reference](https://www.moesif.com/docs/api?go#update-users-in-batch).

```go

import (
	moesifmiddleware "github.com/moesif/moesifmiddleware-go"
)

func literalFieldValue(value string) *string {
    return &value
}

var moesifOptions = map[string]interface{} {
	"Application_Id": "Your Moesif Application Id",
}

// List of Users
var users []*models.UserModel

// Campaign object is optional, but useful if you want to track ROI of acquisition channels
// See https://www.moesif.com/docs/api#users for campaign schema
campaign := models.CampaignModel {
  UtmSource: literalFieldValue("google"),
  UtmMedium: literalFieldValue("cpc"), 
  UtmCampaign: literalFieldValue("adwords"),
  UtmTerm: literalFieldValue("api+tooling"),
  UtmContent: literalFieldValue("landing"),
}
  
// metadata can be any custom dictionary
metadata := map[string]interface{}{
  "email": "john@acmeinc.com",
  "first_name": "John",
  "last_name": "Doe",
  "title": "Software Engineer",
  "sales_info": map[string]interface{}{
      "stage": "Customer",
      "lifetime_value": 24000,
      "account_owner": "mary@contoso.com",
  },
}

// Only UserId is required
userA := models.UserModel{
  UserId:  "12345",
  CompanyId:  literalFieldValue("67890"), // If set, associate user with a company object
  Campaign:  &campaign,
  Metadata:  &metadata,
}

users = append(users, &userA)

// Update User
moesifmiddleware.UpdateUsersBatch(users, moesifOption)
```

## Update Company

### UpdateCompany method
Create or update a company profile in Moesif.
The metadata field can be any company demographic or other info you want to store.
Only the `CompanyId` field is required.
This method is a convenient helper that calls the Moesif API lib.
For details, visit the [Go API Reference](https://www.moesif.com/docs/api?go#update-a-company).

```go
import (
	moesifmiddleware "github.com/moesif/moesifmiddleware-go"
)

func literalFieldValue(value string) *string {
    return &value
}

var moesifOptions = map[string]interface{} {
	"Application_Id": "Your Moesif Application Id",
}

// Campaign object is optional, but useful if you want to track ROI of acquisition channels
// See https://www.moesif.com/docs/api#update-a-company for campaign schema
campaign := models.CampaignModel {
  UtmSource: literalFieldValue("google"),
  UtmMedium: literalFieldValue("cpc"), 
  UtmCampaign: literalFieldValue("adwords"),
  UtmTerm: literalFieldValue("api+tooling"),
  UtmContent: literalFieldValue("landing"),
}
  
// metadata can be any custom dictionary
metadata := map[string]interface{}{
  "org_name": "Acme, Inc",
  "plan_name": "Free",
  "deal_stage": "Lead",
  "mrr": 24000,
  "demographics": map[string]interface{}{
      "alexa_ranking": 500000,
      "employee_count": 47,
  },
}

// Prepare company model
company := models.CompanyModel{
	CompanyId:		  "67890",	// The only required field is your company id
	CompanyDomain:  literalFieldValue("acmeinc.com"), // If domain is set, Moesif will enrich your profiles with publicly available info 
	Campaign: 		  &campaign,
	Metadata:		    &metadata,
}

// Update Company
moesifmiddleware.UpdateCompany(&company, moesifOption)
```

### UpdateCompaniesBatch method
Similar to UpdateCompany, but used to update a list of companies in one batch. 
Only the `CompanyId` field is required.
This method is a convenient helper that calls the Moesif API lib.
For details, visit the [Go API Reference](https://www.moesif.com/docs/api?go#update-companies-in-batch).

```go

import (
	moesifmiddleware "github.com/moesif/moesifmiddleware-go"
)

func literalFieldValue(value string) *string {
    return &value
}

var moesifOptions = map[string]interface{} {
	"Application_Id": "Your Moesif Application Id",
}

// List of Companies
var companies []*models.CompanyModel

// Campaign object is optional, but useful if you want to track ROI of acquisition channels
// See https://www.moesif.com/docs/api#update-a-company for campaign schema
campaign := models.CampaignModel {
  UtmSource: literalFieldValue("google"),
  UtmMedium: literalFieldValue("cpc"), 
  UtmCampaign: literalFieldValue("adwords"),
  UtmTerm: literalFieldValue("api+tooling"),
  UtmContent: literalFieldValue("landing"),
}
  
// metadata can be any custom dictionary
metadata := map[string]interface{}{
  "org_name": "Acme, Inc",
  "plan_name": "Free",
  "deal_stage": "Lead",
  "mrr": 24000,
  "demographics": map[string]interface{}{
      "alexa_ranking": 500000,
      "employee_count": 47,
  },
}

// Prepare company model
companyA := models.CompanyModel{
	CompanyId:		  "67890",	// The only required field is your company id
	CompanyDomain:  literalFieldValue("acmeinc.com"), // If domain is set, Moesif will enrich your profiles with publicly available info 
	Campaign: 		  &campaign,
	Metadata:		    &metadata,
}

companies = append(companies, &companyA)

// Update Companies
moesifmiddleware.UpdateCompaniesBatch(companies, moesifOption)
```

## Example

- An example app with Moesif integration is available __[on GitHub](https://github.com/Moesif/moesifmiddleware-go-example).__
- An example for Google Cloud Platform Serverless: __[Google Cloud Functions](https://github.com/Moesif/moesif-gcp-function-go-example)__

## Other integrations

To view more documentation on integration options, please visit __[the Integration Options Documentation](https://www.moesif.com/docs/getting-started/integration-options/).__

[ico-built-for]: https://img.shields.io/badge/built%20for-go-blue.svg
[ico-license]: https://img.shields.io/badge/License-Apache%202.0-green.svg
[ico-source]: https://img.shields.io/github/last-commit/moesif/moesifmiddleware-go.svg?style=social

[link-built-for]: https://golang.org/
[link-license]: https://raw.githubusercontent.com/Moesif/moesifmiddleware-go/master/LICENSE
[link-source]: https://github.com/Moesif/moesifmiddleware-go
