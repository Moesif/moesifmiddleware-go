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
		if debug {
			log.Printf("Sending the outgoing event to Moesif")
		}
		
		// Check if the event is to Moesif
		if !(strings.Contains(request.URL.String(), "moesif.net")) {
		
			// Get Request Body
			var outgoingReqBody interface{}
			if request.Body != nil {
				copyBody, err := request.GetBody()
				if err != nil {
					if debug{
						log.Printf("Error while getting the request body: %s.\n", err.Error())
					}
				}
			
				// Read the request body
				readReqBody, reqBodyErr := ioutil.ReadAll(copyBody)
				if reqBodyErr != nil {	
					if debug {
						log.Printf("Error while reading request body: %s.\n", reqBodyErr.Error())
					}
				}
			
				// Parse the request Body
				if jsonMarshalErr := json.Unmarshal(readReqBody, &outgoingReqBody); jsonMarshalErr != nil {
					if debug {
						log.Printf("Error while parsing request body: %s.\n", jsonMarshalErr.Error())
					}
					outgoingReqBody = nil
				}
			
				// Return io.ReadCloser while making sure a Close() is available for request body
				request.Body = ioutil.NopCloser(bytes.NewBuffer(readReqBody))
			
			} else {
				// Empty Request body
				outgoingReqBody = nil
			}
		
			// Get Response Body
			var outgoingRespBody interface{}
			if response.Body != nil {
				// Read the response body
				bodyBytes, err := ioutil.ReadAll(response.Body)
				if err != nil {
					if debug {
						log.Printf("Error while reading response body: %s.\n", err.Error())
					}
				}
			 
				// Convert response body into string
				outgoingRespBody = string(bodyBytes)
	
				// Return io.ReadCloser while making sure a Close() is available for response body
				response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			} else {
				// Empty response Body
				outgoingRespBody = nil
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
		
			// Get Outgoing Session Token
			var sessionTokenOutgoing string
			if _, found := moesifOption["Get_Session_Token_Outgoing"]; found {
				sessionTokenOutgoing = moesifOption["Get_Session_Token_Outgoing"].(func(*http.Request, *http.Response) string)(request, response)
			}

			// Send Event To Moesif
			sendMoesifAsync(request, outgoingReqTime, nil, outgoingReqBody, outgoingRspTime, response.StatusCode, 
							response.Header, outgoingRespBody, &userIdOutgoing, &sessionTokenOutgoing, metadataOutgoing)
			
			} else {
				if debug {
					log.Println("Skip sending Moesif Event")
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
