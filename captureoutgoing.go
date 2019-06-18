package moesifmiddleware

import (
	"context"
	"log"
	"net/http"
	"time"
	"strings"
	"io/ioutil"
	"encoding/json"
	"bytes"
	b64 "encoding/base64"
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
		if debug{
			log.Printf("Skip sending the outgoing event to Moesif")
		}
	} else {
		
		// Check if the event is to Moesif
		if !(strings.Contains(request.URL.String(), "moesif.net")) {

			if debug {
				log.Printf("Sending the outgoing event to Moesif")
			}

			// Get Request Body
			var outgoingReqBody interface{}
			var reqEncoding string
			if request.Body != nil {
				copyBody, err := request.GetBody()
				if err != nil {
					if debug{
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
			
				// Parse the request Body
				reqEncoding = "json"
				if jsonReqParseErr := json.Unmarshal(readReqBody, &outgoingReqBody); jsonReqParseErr != nil {
					if debug {
						log.Printf("About to parse outgoing request body as base64 ")
					}
					outgoingReqBody = b64.StdEncoding.EncodeToString(readReqBody)
					reqEncoding = "base64"
					if debug {
						log.Printf("Parsed outgoing request body as base64 - %s", outgoingReqBody)
					}
				}
			
				// Return io.ReadCloser while making sure a Close() is available for request body
				request.Body = ioutil.NopCloser(bytes.NewBuffer(readReqBody))
			
			} else {
				// Empty Request body
				outgoingReqBody = nil
				reqEncoding = ""
			}
  
			// Get Response Body
			var outgoingRespBody interface{}
			var respEncoding string
			if response.Body != nil {
				// Read the response body
				readRespBody, err := ioutil.ReadAll(response.Body)
				if err != nil {
					if debug {
						log.Printf("Error while reading outgoing response body: %s.\n", err.Error())
					}
				}

				// Parse the response Body
				respEncoding = "json"
				if jsonRespParseErr := json.Unmarshal(readRespBody, &outgoingRespBody); jsonRespParseErr != nil {
					if debug {
						log.Printf("About to parse outgoing response body as base64 ")
					}
					// Base64 Encode data
					outgoingRespBody = b64.StdEncoding.EncodeToString(readRespBody)
					respEncoding = "base64"
					if debug {
						log.Printf("Parsed outgoing response body as base64 - %s", outgoingRespBody)
					}
				}
	
				// Return io.ReadCloser while making sure a Close() is available for response body
				response.Body = ioutil.NopCloser(bytes.NewBuffer(readRespBody))
			} else {
				if debug {
					log.Printf("Error while parsing outgoing response body ")
				}
				// Empty response Body
				outgoingRespBody = nil
				respEncoding = ""
			}
			
		
			// Get Outgoing Event Metadata
			var metadataOutgoing map[string]interface{} = nil
			if _, found := moesifOption["Get_Metadata_Outgoing"]; found {
				metadataOutgoing = moesifOption["Get_Metadata_Outgoing"].(func(*http.Request, *http.Response) map[string]interface{})(request, response)
			}
		
			// Get Outgoing User
			var userIdOutgoing string
			if _, found := moesifOption["Identify_User_Outgoing"]; found {
				userIdOutgoing = moesifOption["Identify_User_Outgoing"].(func(*http.Request, *http.Response) string)(request, response)
			}

			// Get Outgoing Company
			var companyIdOutgoing string
			if _, found := moesifOption["Identify_Company_Outgoing"]; found {
				companyIdOutgoing = moesifOption["Identify_Company_Outgoing"].(func(*http.Request, *http.Response) string)(request, response)
			}
		
			// Get Outgoing Session Token
			var sessionTokenOutgoing string
			if _, found := moesifOption["Get_Session_Token_Outgoing"]; found {
				sessionTokenOutgoing = moesifOption["Get_Session_Token_Outgoing"].(func(*http.Request, *http.Response) string)(request, response)
			}

			// Send Event To Moesif
			sendMoesifAsync(request, outgoingReqTime, nil, outgoingReqBody, &reqEncoding, outgoingRspTime, response.StatusCode, 
							response.Header, outgoingRespBody, &respEncoding, &userIdOutgoing, &companyIdOutgoing, &sessionTokenOutgoing, metadataOutgoing)
			
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
