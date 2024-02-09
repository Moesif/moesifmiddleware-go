package moesifmiddleware

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// Transport implements http.RoundTripper.
type Transport struct {
	Transport   http.RoundTripper
	LogRequest  func(req *http.Request)
	LogResponse func(resp *http.Response)
}

// The default logging transport that wraps http.DefaultTransport.
var DefaultTransport = &Transport{
	Transport: http.DefaultTransport,
}

type contextKey struct {
	name string
}

var ContextKeyRequestStart = &contextKey{"RequestStart"}

// RoundTrip is the core part of this module and implements http.RoundTripper.
func (t *Transport) RoundTrip(request *http.Request) (*http.Response, error) {
	ctx := context.WithValue(request.Context(), ContextKeyRequestStart, time.Now())
	request = request.WithContext(ctx)

	// Outgoing Request Time
	outgoingReqTime := time.Now().UTC()

	response, err := t.transport().RoundTrip(request)
	if err != nil {
		return response, err
	}

	// Outgoing Response Time
	outgoingRspTime := time.Now().UTC()

	// Skip capture outgoing event
	shouldSkipOutgoing := false
	if _, found := moesifOption["Should_Skip_Outgoing"]; found {
		shouldSkipOutgoing = moesifOption["Should_Skip_Outgoing"].(func(*http.Request, *http.Response) bool)(request, response)
	}

	// Skip / Send event to moesif
	if shouldSkipOutgoing {
		if debug {
			log.Printf("Skip sending the outgoing event to Moesif")
		}
	} else {

		// Check if the event is to Moesif
		if !(strings.Contains(request.URL.String(), "moesif.net")) {

			if debug {
				log.Printf("Sending the outgoing event to Moesif")
			}

			// Get Request Body
			var (
				outgoingReqBody  interface{}
				reqEncoding      string
				reqContentLength *int64
			)

			if logBodyOutgoing && request.Body != nil {
				copyBody, err := request.GetBody()
				if err != nil {
					if debug {
						log.Printf("Error while getting the outgoing request body: %s.\n", err.Error())
					}
				}

				// Read the request body
				readReqBody, reqBodyErr := ioutil.ReadAll(copyBody)
				if reqBodyErr != nil {
					if debug {
						log.Printf("Error while reading outgoing request body: %s.\n", reqBodyErr.Error())
					}
				}
				reqContentLength = getContentLength(request.Header, readReqBody)

				// Parse the request Body
				outgoingReqBody, reqEncoding = parseBody(readReqBody, "Request_Body_Masks")

				// Return io.ReadCloser while making sure a Close() is available for request body
				request.Body = ioutil.NopCloser(bytes.NewBuffer(readReqBody))

			}

			// Get Response Body
			var (
				outgoingRespBody  interface{}
				respEncoding      string
				respContentLength *int64
			)

			if logBodyOutgoing && response.Body != nil {
				// Read the response body
				readRespBody, err := ioutil.ReadAll(response.Body)
				if err != nil {
					if debug {
						log.Printf("Error while reading outgoing response body: %s.\n", err.Error())
					}
				}
				respContentLength = getContentLength(response.Header, readRespBody)

				// Parse the response Body
				outgoingRespBody, respEncoding = parseBody(readRespBody, "Response_Body_Masks")

				// Return io.ReadCloser while making sure a Close() is available for response body
				response.Body = ioutil.NopCloser(bytes.NewBuffer(readRespBody))
			}

			// Get Outgoing Event Metadata
			var metadataOutgoing map[string]interface{} = nil
			if _, found := moesifOption["Get_Metadata_Outgoing"]; found {
				metadataOutgoing = moesifOption["Get_Metadata_Outgoing"].(func(*http.Request, *http.Response) map[string]interface{})(request, response)
			}

			// Get Outgoing User
			userIdOutgoing := getConfigStringValuesForOutgoingEvent("Identify_User_Outgoing", request, response)

			// Get Outgoing Company
			companyIdOutgoing := getConfigStringValuesForOutgoingEvent("Identify_Company_Outgoing", request, response)

			// Get Outgoing Session Token
			sessionTokenOutgoing := getConfigStringValuesForOutgoingEvent("Get_Session_Token_Outgoing", request, response)

			direction := "Outgoing"

			// Mask Request Header
			var requestHeader map[string]interface{}
			requestHeader = maskHeaders(HeaderToMap(request.Header), "Request_Header_Masks")

			// Mask Response Header
			var responseHeader map[string]interface{}
			responseHeader = maskHeaders(HeaderToMap(response.Header), "Response_Header_Masks")

			// Send Event To Moesif
			sendMoesifAsync(request, outgoingReqTime, requestHeader, nil, outgoingReqBody, &reqEncoding, reqContentLength,
				outgoingRspTime, response.StatusCode, responseHeader, outgoingRespBody, &respEncoding, respContentLength,
				userIdOutgoing, companyIdOutgoing, &sessionTokenOutgoing, metadataOutgoing, &direction)

		} else {
			if debug {
				log.Println("Request Skipped since it is Moesif Event")
			}
		}
	}

	return response, err
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}

	return http.DefaultTransport
}
