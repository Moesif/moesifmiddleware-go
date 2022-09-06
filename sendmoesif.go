package moesifmiddleware

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/moesif/moesifapi-go/models"
)

// Send Event to Moesif
func sendMoesifAsync(request *http.Request, reqTime time.Time, reqHeader map[string]interface{}, apiVersion *string, reqBody interface{}, reqEncoding *string,
	rspTime time.Time, respStatus int, respHeader map[string]interface{}, respBody interface{}, respEncoding *string,
	userId string, companyId string, sessionToken *string, metadata map[string]interface{},
	direction *string) {

	// Get Client Ip
	ip := getClientIp(request)

	// Prepare request model
	event_request := models.EventRequestModel{
		Time:             &reqTime,
		Uri:              request.URL.Scheme + "://" + request.Host + request.URL.Path,
		Verb:             request.Method,
		ApiVersion:       apiVersion,
		IpAddress:        &ip,
		Headers:          reqHeader,
		Body:             &reqBody,
		TransferEncoding: reqEncoding,
	}

	// Prepare response model
	event_response := models.EventResponseModel{
		Time:             &rspTime,
		Status:           respStatus,
		IpAddress:        nil,
		Headers:          respHeader,
		Body:             respBody,
		TransferEncoding: respEncoding,
	}

	// Generate random percentage
	rand.Seed(time.Now().UnixNano())
	randomPercentage := rand.Intn(100)

	// Parse sampling percentage based on user/company
	samplingPercentage := getSamplingPercentage(userId, companyId)

	if samplingPercentage > randomPercentage {

		// Add Weight to the Event Model
		var eventWeight int
		if samplingPercentage == 0 {
			eventWeight = 1
		} else {
			eventWeight = 100 / samplingPercentage
		}

		// Prepare the event model
		event := models.EventModel{
			Request:      event_request,
			Response:     event_response,
			SessionToken: sessionToken,
			Tags:         nil,
			UserId:       &userId,
			CompanyId:    &companyId,
			Metadata:     metadata,
			Direction:    direction,
			Weight:       &eventWeight,
		}

		errSendEvent := apiClient.QueueEvent(&event)
		if errSendEvent != nil {
			log.Fatalf("Error while adding event to Moesif: %s.\n", errSendEvent.Error())
		} else {
			if debug {
				log.Println("Event successfully added to the queue")
			}
		}
	} else {
		if debug{
			log.Println("Skipped Event due to sampling percentage: " + strconv.Itoa(samplingPercentage) + " and random percentage: " + strconv.Itoa(randomPercentage))
		}
	}
}
