/*
 * moesifmiddleware-go
 */
package moesifmiddleware

import (
	"bytes"
	"crypto/rand"
	"fmt"
	moesifapi "github.com/moesif/moesifapi-go"
	"github.com/moesif/moesifapi-go/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Global variable
var (
	apiClient              moesifapi.API
	debug                  bool
	moesifOption           map[string]interface{}
	disableCaptureOutgoing bool
	disableTransactionId   bool
	logBody                bool
	logBodyOutgoing        bool
	samplingPercentage     int
	eTag                   string
	lastUpdatedTime        time.Time
	appConfig              map[string]interface{}
	userSampleRateMap      map[string]interface{}
	companySampleRateMap   map[string]interface{}
)

// Initialize the client
func moesifClient(moesifOption map[string]interface{}) {

	var apiEndpoint string
	var batchSize int
	var eventQueueSize int
	var timerWakeupSeconds int

	// Try to fetch the api endpoint
	if endpoint, found := moesifOption["Api_Endpoint"].(string); found {
		apiEndpoint = endpoint
	}

	// Try to fetch the event queue size
	if queueSize, found := moesifOption["Event_Queue_Size"].(int); found {
		eventQueueSize = queueSize
	}

	// Try to fetch the batch size
	if batch, found := moesifOption["Batch_Size"].(int); found {
		batchSize = batch
	}

	// Try to fetch the timer wake up seconds
	if timer, found := moesifOption["Timer_Wake_Up_Seconds"].(int); found {
		timerWakeupSeconds = timer
	}

	api := moesifapi.NewAPI(moesifOption["Application_Id"].(string), &apiEndpoint, eventQueueSize, batchSize, timerWakeupSeconds)
	apiClient = api

	//  Disable debug by default
	debug = false
	// Try to fetch the debug from the option
	if isDebug, found := moesifOption["Debug"].(bool); found {
		debug = isDebug
	}

	// Disable TransactionId by default
	disableTransactionId = false
	// Try to fetch the disableTransactionId from the option
	if isEnabled, found := moesifOption["disableTransactionId"].(bool); found {
		disableTransactionId = isEnabled
	}

	// Enable logBody by default
	logBody = true
	// Try to fetch the disableTransactionId from the option
	if isEnabled, found := moesifOption["Log_Body"].(bool); found {
		logBody = isEnabled
	}

	// Fetch application config
	response, err := apiClient.GetAppConfig()

	// Parse Configuration
	if err == nil {
		samplingPercentage, eTag, lastUpdatedTime = parseConfiguration(response)
	} else {
		if debug {
			log.Printf("Error fetching application configuration on initilization with err: %s.\n", err.Error())
		}
	}
}

// Moesif Response Recorder
type MoesifResponseRecorder struct {
	rw     http.ResponseWriter
	status int
	writer io.Writer
	header map[string][]string
}

// Function to generate UUID
func uuid() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

// Response Recorder
func responseRecorder(rw http.ResponseWriter, status int, writer io.Writer) MoesifResponseRecorder {
	rr := MoesifResponseRecorder{
		rw,
		status,
		writer,
		make(map[string][]string, 5),
	}
	return rr
}

// Implementing the WriteHeader method of ResponseWriter Interface
func (rec *MoesifResponseRecorder) WriteHeader(code int) {
	rec.status = code
	rec.rw.WriteHeader(code)
}

// Implementing the Write method of ResponseWriter Interface
func (rec *MoesifResponseRecorder) Write(b []byte) (int, error) {
	return rec.writer.Write(b)
}

// Implementing the Header method of ResponseWriter Interface
func (rec *MoesifResponseRecorder) Header() http.Header {
	return rec.rw.Header()
}

// Start Capture Outgoing Request
func StartCaptureOutgoing(configurationOption map[string]interface{}) {

	// Call the function to initialize the moesif client and moesif options
	if apiClient == nil {
		// Set the Capture_Outoing_Requests to true to capture outgoing request
		configurationOption["Capture_Outoing_Requests"] = true
		moesifOption = configurationOption
		moesifClient(moesifOption)
	}

	if debug {
		log.Println("Start Capturing outgoing requests")
	}
	// Enable logBody by default
	logBodyOutgoing = true
	// Try to fetch the disableTransactionId from the option
	if isEnabled, found := moesifOption["Log_Body_Outgoing"].(bool); found {
		logBodyOutgoing = isEnabled
	}

	http.DefaultTransport = DefaultTransport
}

// Update User
func UpdateUser(user *models.UserModel, configurationOption map[string]interface{}) {

	// Call the function to initialize the moesif client and moesif options
	if apiClient == nil {
		moesifClient(configurationOption)
	}

	// Add event to the queue
	errUpdateUser := apiClient.QueueUser(user)
	// Log the message
	if errUpdateUser != nil {
		log.Fatalf("Error while updating user: %s.\n", errUpdateUser.Error())
	} else {
		log.Println("Update User successfully added to the queue")
	}
}

// Update Users Batch
func UpdateUsersBatch(users []*models.UserModel, configurationOption map[string]interface{}) {

	// Call the function to initialize the moesif client and moesif options
	if apiClient == nil {
		moesifClient(configurationOption)
	}

	// Add event to the queue
	errUpdateUserBatch := apiClient.QueueUsers(users)
	// Log the message
	if errUpdateUserBatch != nil {
		log.Fatalf("Error while updating users in batch: %s.\n", errUpdateUserBatch.Error())
	} else {
		log.Println("Updated Users successfully added to the queue")
	}
}

// Update Company
func UpdateCompany(company *models.CompanyModel, configurationOption map[string]interface{}) {

	// Call the function to initialize the moesif client and moesif options
	if apiClient == nil {
		moesifClient(configurationOption)
	}

	// Add event to the queue
	errUpdateCompany := apiClient.QueueCompany(company)
	// Log the message
	if errUpdateCompany != nil {
		log.Fatalf("Error while updating company: %s.\n", errUpdateCompany.Error())
	} else {
		log.Println("Update Company successfully added to the queue")
	}
}

// Update Companies Batch
func UpdateCompaniesBatch(companies []*models.CompanyModel, configurationOption map[string]interface{}) {

	// Call the function to initialize the moesif client and moesif options
	if apiClient == nil {
		moesifClient(configurationOption)
	}

	// Add event to the queue
	errUpdateCompaniesBatch := apiClient.QueueCompanies(companies)
	// Log the message
	if errUpdateCompaniesBatch != nil {
		log.Fatalf("Error while updating companies in batch: %s.\n", errUpdateCompaniesBatch.Error())
	} else {
		log.Println("Updated companies successfully added to the queue")
	}
}

// Moesif Middleware
func MoesifMiddleware(next http.Handler, configurationOption map[string]interface{}) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		// Buffer
		var buf bytes.Buffer

		// Create a writer to duplicates it's writes to all the provided writers
		multiWriter := io.MultiWriter(rw, &buf)

		// Initialize the status to 200 in case WriteHeader is not called
		response := responseRecorder(
			rw,
			200,
			multiWriter,
		)

		// Call the function to initialize the moesif client and moesif options
		if apiClient == nil {
			moesifOption = configurationOption
			moesifClient(moesifOption)
		}

		// Add transactionId to the headers
		if !disableTransactionId {
			// Try to fetch the transactionId from the header
			transactionId := request.Header.Get("X-Moesif-Transaction-Id")
			// Check if need to generate transactionId
			if len(transactionId) == 0 {
				transactionId, _ = uuid()
			}

			if len(transactionId) != 0 {
				// Add transationId to the request header
				request.Header.Set("X-Moesif-Transaction-Id", transactionId)

				// Add transationId to the response header
				rw.Header().Add("X-Moesif-Transaction-Id", transactionId)
			}
		}

		// Request Time
		requestTime := time.Now().UTC()

		// Serve the HTTP Request
		next.ServeHTTP(&response, request)

		// Response Time
		responseTime := time.Now().UTC()

		shouldSkip := false
		if _, found := moesifOption["Should_Skip"]; found {
			shouldSkip = moesifOption["Should_Skip"].(func(*http.Request, MoesifResponseRecorder) bool)(request, response)
		}

		if shouldSkip {
			if debug {
				log.Printf("Skip sending the event to Moesif")
			}
		} else {
			if debug {
				log.Printf("Sending the event to Moesif")
			}
			// Call the function to send event to Moesif
			sendEvent(request, response, buf.String(), requestTime, responseTime)
		}
	})
}

// Sending event to Moesif
func sendEvent(request *http.Request, response MoesifResponseRecorder, rspBufferString string, reqTime time.Time, rspTime time.Time) {
	// Get Api Version
	var apiVersion *string = nil
	if isApiVersion, found := moesifOption["Api_Version"].(string); found {
		apiVersion = &isApiVersion
	}

	// Get Request Body
	var reqBody interface{}
	var reqEncoding string
	readReqBody, reqBodyErr := ioutil.ReadAll(request.Body)
	if reqBodyErr != nil {
		if debug {
			log.Printf("Error while reading request body: %s.\n", reqBodyErr.Error())
		}
	}

	// Check if the request body is empty
	reqBody = nil
	if logBody && (len(readReqBody)) > 0 {
		reqBody, reqEncoding = parseBody(readReqBody, "Request_Body_Masks")
	}

	// Get the response body
	var respBody interface{}
	var respEncoding string

	// Parse the response Body
	respBody = nil
	if logBody {
		respBody, respEncoding = parseBody([]byte(rspBufferString), "Response_Body_Masks")
	}

	// Get URL Scheme
	if request.URL.Scheme == "" {
		request.URL.Scheme = "http"
	}

	// Get Metadata
	var metadata map[string]interface{} = nil
	if _, found := moesifOption["Get_Metadata"]; found {
		metadata = moesifOption["Get_Metadata"].(func(*http.Request, MoesifResponseRecorder) map[string]interface{})(request, response)
	}

	// Get User
	userId := getConfigStringValuesForIncomingEvent("Identify_User", request, response)

	// Get Company
	companyId := getConfigStringValuesForIncomingEvent("Identify_Company", request, response)

	// Get Session Token
	sessionToken := getConfigStringValuesForIncomingEvent("Get_Session_Token", request, response)

	direction := "Incoming"

	// Mask Request Header
	var requestHeader map[string]interface{}
	requestHeader = maskHeaders(HeaderToMap(request.Header), "Request_Header_Masks")

	// Mask Response Header
	var responseHeader map[string]interface{}
	responseHeader = maskHeaders(HeaderToMap(response.Header()), "Response_Header_Masks")

	// Send Event To Moesif
	sendMoesifAsync(request, reqTime, requestHeader, apiVersion, reqBody, &reqEncoding, rspTime, response.status, responseHeader, respBody, &respEncoding, userId, companyId, &sessionToken, metadata, &direction)
}
