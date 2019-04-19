package moesifmiddleware

import (
	"log"
	"time"
	"net/http"
	"github.com/moesif/moesifapi-go/models"
)

// Send Event to Moesif
func sendMoesifAsync(request *http.Request, reqTime time.Time, apiVersion *string, reqBody interface{}, 
					 rspTime time.Time, respStatus int, respHeader http.Header, respBody interface{},
					 userId *string, sessionToken *string, metadata map[string]interface{}) {
	
	// Prepare request model
	event_request := models.EventRequestModel{
		Time:       &reqTime,
		Uri:        request.URL.Scheme + "://" + request.Host + request.URL.Path,
		Verb:       request.Method,
		ApiVersion: apiVersion,
		IpAddress:  nil,
		Headers:    request.Header,
		Body: 		&reqBody,
	}

	// Prepare response model
	event_response := models.EventResponseModel{
		Time:      &rspTime,
		Status:    respStatus,
		IpAddress: nil,
		Headers:   respHeader,
		Body: 	   respBody,
	}
	
	// Prepare the event model
	event := models.EventModel{
		Request:      event_request,
		Response:     event_response,
		SessionToken: sessionToken,
		Tags:         nil,
		UserId:       userId,
		CompanyId:    nil,
		Metadata: 	  metadata,
	}

	// Add event to the queue
	errSendEvent := apiClient.QueueEvent(&event)
	// Log the message
	if errSendEvent != nil {
		log.Fatalf("Error while adding event to Moesif: %s.\n", errSendEvent.Error())
	} else {
		log.Println("Event successfully added to the queue")
	}
}
