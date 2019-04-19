# Moesif Middleware for Go

[![Built For][ico-built-for]][link-built-for]
[![Software License][ico-license]][link-license]
[![Source Code][ico-source]][link-source]

Send API call data to [Moesif API Analytics](https://www.moesif.com/)

[Source Code on GitHub](https://github.com/moesif/moesifmiddleware-go)

## How to install
Run the following commands:

```shell
go get github.com/moesif/moesifmiddleware-go
```

## Configuration options

Please note the request is the original http.Request and the response is moesifResponseRecorder.

```go
// Moesif Response Recorder
 type MoesifResponseRecorder struct {
	 rw http.ResponseWriter
	 status int
	 writer io.Writer
	 header map[string][]string
 }
```

#### __`Application_Id`__
(__required__), _string_, is obtained via your Moesif Account, this is required.

#### __`Should_Skip`__
(optional) _(request, response) => boolean_, a function that takes a request and a response,
and returns true if you want to skip this particular event.

#### __`Identify_User`__
(optional, but highly recommended) _(request, response) => string_, a function that takes a request and response, and returns a string that is the user id used by your system. While Moesif tries to identify users automatically, but different frameworks and your implementation might be very different, it would be helpful and much more accurate to provide this function.

#### __`Get_Metadata`__
(optional) _(request, response) => dictionary_, a function that takes a request and response, and
returns a dictionary (must be able to be encoded into JSON). This allows you
to associate this event with custom metadata. For example, you may want to save a VM instance_id, a trace_id, or a tenant_id with the request.

#### __`Get_Session_Token`__
(optional) _(request, response) => string_, a function that takes a request and response, and returns a string that is the session token for this event. Moesif tries to get the session token automatically, but if this doesn't work for your service, you should use this to identify sessions.

#### __`Mask_Event_Model`__
(optional) _(EventModel) => EventModel_, a function that takes an EventModel and returns an EventModel with desired data removed. The return value must be a valid EventModel required by Moesif data ingestion API. For details regarding EventModel please see the [Moesif Golang API Documentation](https://www.moesif.com/docs/api?go).

#### __`Debug`__
(optional) _boolean_, a flag to see debugging messages.

#### __`Capture_Outoing_Requests`__
(optional) _boolean_, Default False. Set to True to capture all outgoing API calls from your app to third parties like Stripe or to your own dependencies while using [net/http](https://golang.org/pkg/net/http/) package. The options below is applied to outgoing API calls.
When the request is outgoing, for options functions that take request and response as input arguments, the request and response objects passed in are [Request](https://golang.org/src/net/http/request.go) request and [Response](https://golang.org/src/net/http/response.go) response objects.

##### __`Should_Skip_Outgoing`__
(optional) _(request, response) => boolean_, a function that takes a request and response, and returns true if you want to skip this particular event.

##### __`Identify_User_Outgoing`__
(optional, but highly recommended) _(request, response) => string_, a function that takes request and response, and returns a string that is the user id used by your system. While Moesif tries to identify users automatically,
but different frameworks and your implementation might be very different, it would be helpful and much more accurate to provide this function.

##### __`Get_Metadata_Outgoing`__
(optional) _(request, response) => dictionary_, a function that takes request and response, and
returns a dictionary (must be able to be encoded into JSON). This allows
to associate this event with custom metadata. For example, you may want to save a VM instance_id, a trace_id, or a tenant_id with the request.

##### __`Get_Session_Token_Outgoing`__
(optional) _(request, response) => string_, a function that takes request and response, and returns a string that is the session token for this event. Again, Moesif tries to get the session token automatically, but if you setup is very different from standard, this function will be very help for tying events together, and help you replay the events.

##### __`Mask_Event_Model_Outgoing`__
(optional) _(EventModel) => EventModel_, a function that takes an EventModel and returns an EventModel with desired data removed. The return value must be a valid EventModel required by Moesif data ingestion API. For details regarding EventModel please see the [Moesif Golang API Documentation](https://www.moesif.com/docs/api?go).

## Example
An example app with Moesif integration is available __[on GitHub](https://github.com/Moesif/moesifmiddleware-go-example).__

## Other integrations

To view more more documentation on integration options, please visit __[the Integration Options Documentation](https://www.moesif.com/docs/getting-started/integration-options/).__

[ico-built-for]: https://img.shields.io/badge/built%20for-go-blue.svg
[ico-license]: https://img.shields.io/badge/License-Apache%202.0-green.svg
[ico-source]: https://img.shields.io/github/last-commit/moesif/moesifmiddleware-go.svg?style=social

[link-built-for]: https://golang.org/
[link-license]: https://raw.githubusercontent.com/Moesif/moesifmiddleware-go/master/LICENSE
[link-source]: https://github.com/Moesif/moesifmiddleware-go
