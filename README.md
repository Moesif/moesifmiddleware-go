# Moesif Middleware for Go
by [Moesif](https://moesif.com), the [API analytics](https://www.moesif.com/features/api-analytics) and [API monetization](https://www.moesif.com/solutions/metered-api-billing) platform.

[![Built For][ico-built-for]][link-built-for]
[![Software License][ico-license]][link-license]
[![Source Code][ico-source]][link-source]

Moesif middleware for Go logs API calls and sends to [Moesif](https://www.moesif.com) for API analytics and log analysis. This middleware allows you to integrate Moesif's API analytics and 
API monetization features into your Go applications with minimal configuration. 

> If you're new to Moesif, see [our Getting Started](https://www.moesif.com/docs/) resources to quickly get up and running.

## Prerequisites
Before using this middleware, make sure you have the following:

- [An active Moesif account](https://moesif.com/wrap)
- [A Moesif Application ID](#get-your-moesif-application-id)

### Get Your Moesif Application ID
After you log into [Moesif Portal](https://www.moesif.com/wrap), you can get your Moesif Application ID during the onboarding steps. You can always access the Application ID any time by following these steps from Moesif Portal after logging in:

1. Select the account icon to bring up the settings menu.
2. Select **Installation** or **API Keys**.
3. Copy your Moesif Application ID from the **Collector Application ID** field.

<img class="lazyload blur-up" src="images/app_id.png" width="700" alt="Accessing the settings menu in Moesif Portal">

## Install the Middleware
Use `go get`:

```bash
go get github.com/moesif/moesifmiddleware-go
```

If you are using [Go modules](https://go.dev/ref/mod), you can specify a version number as well:

```bash
go get github.com/moesif/moesifmiddleware-go@v1.2.3
```

## Configure the Middleware
See the available [configuration options](#configuration-options) to learn how to configure the middleware for your use case.

## How to Use

The following snippet shows how to use the middleware:

```go
import(
    moesifmiddleware "github.com/moesif/moesifmiddleware-go"
)

func handle(w http.ResponseWriter, r *http.Request) {
	// Your API Logic
}

var moesifOptions = map[string]interface{} {
        "Application_Id": "YOUR_MOESIF_APPLICATION_ID",
        "Log_Body": true,
}
http.Handle("/api", moesifmiddleware.MoesifMiddleware(http.HandlerFunc(handle), moesifOption))
```

Replace *`YOUR_MOESIF_APPLICATION_ID`* with [your Moesif Application ID](#get-your-moesif-application-id).

### Optional: Capturing Outgoing API Calls
In addition to your own APIs, you can also start capturing calls out to third party services through the following method:

```go
moesifmiddleware.StartCaptureOutgoing(moesifOption)
```

#### `handler func(ResponseWriter, *Request)` (Required)

The `handler` function registers the handler function for the given pattern through the `HandlerFunc` adapter. See the [example application code](https://github.com/Moesif/moesifmiddleware-go-example/blob/f3692a169ee0c7e73f109a54f65e28b55c611d01/main.go#L54) for better understanding.

#### `moesifOption` (Required)
A `map[string]interface{}` type containing the configuration options for your application. See [the example application code](https://github.com/Moesif/moesifmiddleware-go-example/blob/f3692a169ee0c7e73f109a54f65e28b55c611d01/moesif_options/moesif_options.go#L111) for better understanding.

See [Configuration Options](#configuration-options) for the common configuration options. See [Options for Logging Outgoing Calls](#options-for-logging-outgoing-calls) for configuration options specific to capturing and logging outgoing API calls.

## Troubleshoot
For a general troubleshooting guide that can help you solve common problems, see [Server Troubleshooting Guide](https://www.moesif.com/docs/troubleshooting/server-troubleshooting-guide/). 

Other troubleshooting supports:

- [FAQ](https://www.moesif.com/docs/faq/)
- [Moesif support email](mailto:support@moesif.com)

## Configuration Options
The following sections describe the available configuration options for this middleware. You can set these options in the Moesif initialization options object. See the the [example application code](https://github.com/Moesif/moesifmiddleware-go-example/blob/master/moesif_options/moesif_options.go) to understand how you can specify these options.

### `Application_Id` (Required)
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
  </tr>
  <tr>
   <td>
    <code>string</code>
   </td>
  </tr>
</table>

A string that [identifies your application in Moesif](#get-your-moesif-application-id).

### `Should_Skip`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Parameters
   </th>
   <th scope="col">
    Return type
   </th>
  </tr>
  <tr>
   <td>
    Function
   </td>
   <td>
    <code>(request, response)</code>
   </td>
   <td>
    <code>boolean</code>
   </td>
  </tr>
</table>

Optional.

A function that takes a request and a response,
and returns `true` if you want to skip this particular event.

### `Identify_User`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Parameters
   </th>
   <th scope="col">
    Return type
   </th>
  </tr>
  <tr>
   <td>
    Function
   </td>
   <td>
    <code>(request, response)</code>
   </td>
   <td>
    <code>string</code>
   </td>
  </tr>
</table>

Optional, but highly recommended.

A function that takes a request and a response, and returns a string that represents the user ID used by your system. 

Moesif identifies users automatically. However, due to the differences arising from different frameworks and implementations, provide this function to ensure user identification properly.

### `Identify_Company`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Parameters
   </th>
   <th scope="col">
    Return type
   </th>
  </tr>
  <tr>
   <td>
    Function
   </td>
   <td>
    <code>(request, response)</code>
   </td>
   <td>
    <code>string</code>
   </td>
  </tr>
</table>

Optional. 

A function that takes a request and response, and returns a string that represents the company ID for this event.

### `Get_Metadata`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Parameters
   </th>
   <th scope="col">
    Return type
   </th>
  </tr>
  <tr>
   <td>
    Function
   </td>
   <td>
    <code>(request, response)</code>
   </td>
   <td>
    <code>dictionary</code>
   </td>
  </tr>
</table>

Optional.

A function that returns an object that allows you to add custom metadata that will be associated with the event. 

The metadata must be a dictionary that can be converted to JSON. For example, you may want to save a virtual machine instance ID, a trace ID, or a tenant ID with the request.

### `Get_Session_Token`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Parameters
   </th>
   <th scope="col">
    Return type
   </th>
  </tr>
  <tr>
   <td>
    Function
   </td>
   <td>
    <code>(request, response)</code>
   </td>
   <td>
    <code>string</code>
   </td>
  </tr>
</table>

Optional.

A function that takes a request and response, and returns a string that represents the session token for this event. 

Similar to users and companies, Moesif tries to retrieve session tokens automatically. But if it doesn't work for your service, provide this function to help identify sessions.

### `Request_Header_Masks`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Parameters
   </th>
   <th scope="col">
    Return type
   </th>
  </tr>
  <tr>
   <td>
    Function
   </td>
   <td>
    <code>()</code>
   </td>
   <td>
    <code>[]string</code>
   </td>
  </tr>
</table>

Optional.

A function that returns an array of strings to mask specific request header fields.

### `Request_Body_Masks`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Parameters
   </th>
   <th scope="col">
    Return type
   </th>
  </tr>
  <tr>
   <td>
    Function
   </td>
   <td>
    <code>()</code>
   </td>
   <td>
    <code>[]string</code>
   </td>
  </tr>
</table>

Optional.

A function that returns array of strings to mask specific request body fields.

### `Response_Header_Masks`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Parameters
   </th>
   <th scope="col">
    Return type
   </th>
  </tr>
  <tr>
   <td>
    Function
   </td>
   <td>
    <code>()</code>
   </td>
   <td>
    <code>[]string</code>
   </td>
  </tr>
</table>

Optional.

A function that returns array of strings to mask specific response header fields.

### `Response_Body_Masks`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Parameters
   </th>
   <th scope="col">
    Return type
   </th>
  </tr>
  <tr>
   <td>
    Function
   </td>
   <td>
    <code>()</code>
   </td>
   <td>
    <code>[]string</code>
   </td>
  </tr>
</table>

Optional.

A function that returns array of strings to mask specific response body fields.

### `Debug`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
  </tr>
  <tr>
   <td>
    <code>boolean</code>
   </td>
  </tr>
</table>

Optional.

Set to `true` to see debugging messages. This may help you troubleshoot integration issues.

### `Log_Body`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Default
   </th>
  </tr>
  <tr>
   <td>
    <code>boolean</code>
   </td>
   <td>
    <code>true</code>
   </td>
  </tr>
</table>

Optional.

Set to `false` to not log the request and response body to Moesif.

### `Event_Queue_Size`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Default
   </th>
  </tr>
  <tr>
   <td>
    <code>int</code>
   </td>
   <td>
    <code>10000</code>
   </td>
  </tr>
</table>

An optional field name that specifies the maximum number of events to hold in queue before sending to Moesif. In case of network issues, the middleware may fail to connect to or send events to Moesif. For those scenarios, this option helps prevent adding new events to the queue to prevent memory overflow.

### `Batch_Size` 
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Default
   </th>
  </tr>
  <tr>
   <td>
    <code>int</code>
   </td>
   <td>
    <code>200</code>
   </td>
  </tr>
</table>

An optional field name that specifies the maximum batch size when sending to Moesif.

### `Timer_Wake_Up_Seconds`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Default
   </th>
  </tr>
  <tr>
   <td>
    <code>int</code>
   </td>
   <td>
    <code>2</code>
   </td>
  </tr>
</table>

An optional field that specifies a time in seconds how often background thread runs to send events to Moesif.

### Options for Logging Outgoing Calls

The following configuration options apply to outgoing API calls. The request and response objects passed in are [`Request`](https://golang.org/src/net/http/request.go) and [`Response`](https://golang.org/src/net/http/response.go) objects of the Go standard library.

### `Should_Skip_Outgoing`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Parameters
   </th>
   <th scope="col">
    Return type
   </th>
  </tr>
  <tr>
   <td>
    Function
   </td>
   <td>
    <code>(request, response)</code>
   </td>
   <td>
    <code>boolean</code>
   </td>
  </tr>
</table>

Optional.

A function that takes a request and response, and returns `true` if you want to skip this particular event.

### `Identify_User_Outgoing`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Parameters
   </th>
   <th scope="col">
    Return type
   </th>
  </tr>
  <tr>
   <td>
    Function
   </td>
   <td>
    <code>(request, response)</code>
   </td>
   <td>
    <code>string</code>
   </td>
  </tr>
</table>

Optional, but highly recommended.

A function that takes a request and a response, and returns a string that represents the user ID used by your system. 

Moesif identifies users automatically. However, due to the differences arising from different frameworks and implementations, provide this function to ensure user identification properly.

### `Identify_Company_Outgoing`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Parameters
   </th>
   <th scope="col">
    Return type
   </th>
  </tr>
  <tr>
   <td>
    Function
   </td>
   <td>
    <code>(request, response)</code>
   </td>
   <td>
    <code>string</code>
   </td>
  </tr>
</table>

Optional.

A function that takes request and response, and returns a string that represents the company ID for this event.

### `Get_Metadata_Outgoing`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Parameters
   </th>
   <th scope="col">
    Return type
   </th>
  </tr>
  <tr>
   <td>
    Function
   </td>
   <td>
    <code>(request, response)</code>
   </td>
   <td>
    <code>dictionary</code>
   </td>
  </tr>
</table>

Optional.

A function that returns an object that allows you to add custom metadata that will be associated with the event. 

The metadata must be a dictionary that can be converted to JSON. For example, you may want to save a virtual machine instance ID, a trace ID, or a tenant ID with the request.

### `Get_Session_Token_Outgoing`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Parameters
   </th>
   <th scope="col">
    Return type
   </th>
  </tr>
  <tr>
   <td>
    Function
   </td>
   <td>
    <code>(request, response)</code>
   </td>
   <td>
    <code>string</code>
   </td>
  </tr>
</table>

Optional.

A function that takes a request and response, and returns a string that represents the session token for this event. 

Similar to users and companies, Moesif tries to retrieve session tokens automatically. But if it doesn't work for your service, provide this function to help identify sessions and replay them.

### `Log_Body_Outgoing`
<table>
  <tr>
   <th scope="col">
    Data type
   </th>
   <th scope="col">
    Default
   </th>
  </tr>
  <tr>
   <td>
    <code>boolean</code>
   </td>
   <td>
    <code>true</code>
   </td>
  </tr>
</table>

Optional.

Set to `false` to not log the request and response body to Moesif.

## Examples

- [Example Go app that using this middleware](https://github.com/Moesif/moesifmiddleware-go-example)
- [Example Go app using this middleware and Google Cloud Run functions](https://github.com/Moesif/moesif-gcp-function-go-example)

The following examples demonstrate some common operations:

- [Updating a single user](#updateuser-method)
- [Updating users in batch](#updateusersbatch-method)
- [Updating a single company](#updatecompany-method)
- [Updating companies in batch](#updatecompaniesbatch-method)
- [Updating a single subscription](#updatesubscription-method)
- [Updating subscriptions in batch](#updatesubscriptionsbatch-method)

## Update User

### `UpdateUser` Method
Use this method to create or update a user profile in Moesif.


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

The `metadata` field can contain any user demographic or other information you want to store.

Only the `UserId` field is required.
This method is a convenient helper that calls the Moesif API library. For more information, see [Moesif Go API documentation](https://www.moesif.com/docs/api?go#update-a-user).

### `UpdateUsersBatch` Method
Similar to `UpdateUser`, but to update a list of users in one batch. 

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

The `metadata` field can contain any company demographic or other information you want to store.

Only the `UserId` field is required.
This method is a convenient helper that calls the Moesif API library. For more information, see [Moesif Go API documentation](https://www.moesif.com/docs/api?go#update-users-in-batch)

## Update Company

### `UpdateCompany` Method
Use this method to create or update a company profile in Moesif.

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

The metadata field can be any company demographic or other info you want to store.

Only the `CompanyId` field is required.

This method is a convenient helper that calls the Moesif API library. For details, see [Moesif Go API documentation](https://www.moesif.com/docs/api?go#update-a-company).


### `UpdateCompaniesBatch` Method
Similar to `UpdateCompany`, but to update a list of companies in one batch. 

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

The metadata field can be any company demographic or other info you want to store.

Only the `CompanyId` field is required.

This method is a convenient helper that calls the Moesif API library. For details, see [Moesif Go API documentation](https://www.moesif.com/docs/api?go#update-companies-in-batch).

## Update Subscription

### `UpdateSubscription` Method
Use this method to create or update a subscription profile in Moesif.
 
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

// metadata can be any custom dictionary
metadata := map[string]interface{}{
  "plan_name": "Pro",
  "deal_stage": "Customer",
  "mrr": 48000,
  "demographics": map[string]interface{}{
      "subscription_length": 12,
      "subscription_type": "annual",
  },
}

// Prepare subscription model
subscription := models.SubscriptionModel{
	SubscriptionId: "12345",	// Required subscription id
  CompanyId: "67890",       // Required company id
	Metadata: &metadata,
}

// Update Subscription
moesifmiddleware.UpdateSubscription(&subscription, moesifOptions)
```


The `metadata` field can be any subscription demographic or other information you want to store.

Only the `SubscriptionId` and `CompanyId` fields are required.

This method is a convenient helper that calls the Moesif API library. For more information, see [Moesif Go API documentation](https://www.moesif.com/docs/api?go#update-a-subscription).


### `UpdateSubscriptionsBatch` method
Similar to `UpdateSubscription`, but to update a list of subscriptions in one batch. 

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

// List of Subscriptions
var subscriptions []*models.SubscriptionModel

// metadata can be any custom dictionary
metadata := map[string]interface{}{
  "plan_name": "Pro",
  "deal_stage": "Customer",
  "mrr": 48000,
  "demographics": map[string]interface{}{
      "subscription_length": 12,
      "subscription_type": "annual",
  },
}

// Prepare subscription model
subscriptionA := models.SubscriptionModel{
	SubscriptionId: "12345",	// Required subscription id
  CompanyId: "67890",       // Required company id
	Metadata: &metadata,
}

subscriptions = append(subscriptions, &subscriptionA)

// Update Subscriptions
moesifmiddleware.UpdateSubscriptionsBatch(subscriptions, moesifOptions)
```

The `metadata` field can be any subscription demographic or other information you want to store.

Only the `SubscriptionId` and `CompanyId` fields are required.

This method is a convenient helper that calls the Moesif API library. For more information, see [Moesif Go API documentation](https://www.moesif.com/docs/api?go#update-subscriptions-in-batch).

## Explore Other Integrations

Explore other integration options from Moesif:

- [Server integration options documentation](https://www.moesif.com/docs/server-integration//)
- [Client integration options documentation](https://www.moesif.com/docs/client-integration/)

[ico-built-for]: https://img.shields.io/badge/built%20for-go-blue.svg
[ico-license]: https://img.shields.io/badge/License-Apache%202.0-green.svg
[ico-source]: https://img.shields.io/github/last-commit/moesif/moesifmiddleware-go.svg?style=social

[link-built-for]: https://golang.org/
[link-license]: https://raw.githubusercontent.com/Moesif/moesifmiddleware-go/master/LICENSE
[link-source]: https://github.com/Moesif/moesifmiddleware-go
