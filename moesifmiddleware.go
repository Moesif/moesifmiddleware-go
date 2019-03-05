/*
 * moesifmiddleware-go
 */
 package moesifmiddleware

 import (
	 "log"
	 "net/http"
	 "bytes"
	 "io"
	 moesifapi "github.com/moesif/moesifapi-go"
	 "github.com/moesif/moesifapi-go/models"
	 "time"
	 "fmt"
	 "crypto/rand"
 )
 
 // Global variable
 var (
	 apiClient moesifapi.API
 )
 
 // Initialize the client
 func moesifClient(moesifOption map[string]interface{}) {
	 api := moesifapi.NewAPI(moesifOption["Application_Id"].(string))
	 apiClient = api
 }
 
 // Moesif Response Recorder
 type moesifResponseRecorder struct {
	 rw http.ResponseWriter
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
 func responseRecorder(rw http.ResponseWriter, status int, writer io.Writer) moesifResponseRecorder{
	 rr := moesifResponseRecorder{
		 rw,
		 status,
		 writer,
		 make(map[string][]string, 5),
	 }
	 return rr
 }
 
 // Implementing the WriteHeader method of ResponseWriter Interface
 func (rec *moesifResponseRecorder) WriteHeader(code int) {
	 rec.status = code
	 rec.rw.WriteHeader(code)
 }
 
 // Implementing the Write method of ResponseWriter Interface
 func (rec *moesifResponseRecorder) Write(b []byte) (int, error){
	 return rec.writer.Write(b)
 }
 
 // Implementing the Header method of ResponseWriter Interface
 func (rec *moesifResponseRecorder) Header() http.Header{
	 return rec.rw.Header()
 }
 
 // Moesif Middleware
 func MoesifMiddleware(next http.Handler, moesifOption map[string]interface{}) http.Handler {
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

		// Disable TransactionId by default
		disableTransactionId := false
		// Try to fetch the disableTransactionId from the option
		if isEnabled, found := moesifOption["disableTransactionId"].(bool); found {
			disableTransactionId = isEnabled
		}

		// Add transactionId to the headers
		if !disableTransactionId {
			// Try to fetch the transactionId from the header
			reqTransId := request.Header.Get("X-Moesif-Transaction-Id")
			// Set the transationId
			transationId := reqTransId
			// Check if need to generate transactionId
			if len(transationId) == 0 {
				transationId, _ = uuid()
			}
			
			if len(transationId) != 0 {
				// Add transationId to the request header
				request.Header.Set("X-Moesif-Transaction-Id", transationId)

				// Add transationId to the response header
				rw.Header().Add("X-Moesif-Transaction-Id", transationId)
			}
		}

		 // Request Time
		 requestTime := time.Now().UTC()

		 // Serve the HTTP Request
		 next.ServeHTTP(&response, request)

		 // Response Time
		 responseTime := time.Now().UTC()
 
		 // Call the function to initialize the moesif client
		 if apiClient == nil {
			moesifClient(moesifOption)
		 }
 
		 // Call the function to send event to Moesif
		 sendEvent(request, response, buf.String(), requestTime, responseTime)
	 })
 }
 
 // Sending event to Moesif
 func sendEvent(request *http.Request, response moesifResponseRecorder, rspBody string, reqTime time.Time, rspTime time.Time) {
 
	 // Prepare request model
	 event_request := models.EventRequestModel{
		 Time:       &reqTime,
		 Uri:        request.Host + request.RequestURI,
		 Verb:       request.Method,
		 ApiVersion: nil,
		 IpAddress:  nil,
		 Headers:    request.Header,
		 Body: 		 nil,
		 }

	// Prepare response model
	event_response := models.EventResponseModel{
		Time:      &rspTime,
		Status:    response.status,
		IpAddress: nil,
		Headers:   response.Header(),
		Body: 	   rspBody,
	}
	 
	 // Prepare the event model
	 event := models.EventModel{
		 Request:      event_request,
		 Response:     event_response,
		 SessionToken: nil,
		 Tags:         nil,
		 UserId:       nil,
		 Metadata: 	   nil,
	 }
 
	 // Add event to the queue
	 err := apiClient.QueueEvent(&event)
	 
	 // Log the message
	 if err != nil {
		 log.Fatalf("Error while adding event to Moesif: %s.\n", err.Error())
	 } else {
		 log.Println("Event successfully added to the queue")
	 }
 }
