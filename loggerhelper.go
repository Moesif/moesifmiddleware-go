package moesifmiddleware

import (
	b64 "encoding/base64"
	"encoding/json"
	"log"
	"net/http"
)

func contains(arr []string, str string) bool {
	for _, value := range arr {
		if value == str {
			return true
		}
	}
	return false
}

func getConfigStringValuesForIncomingEvent(fieldName string, request *http.Request, response MoesifResponseRecorder) string {
	var field string
	if _, found := moesifOption[fieldName]; found {
		field = moesifOption[fieldName].(func(*http.Request, MoesifResponseRecorder) string)(request, response)
	}
	return field
}

func getConfigStringValuesForOutgoingEvent(fieldName string, request *http.Request, response *http.Response) string {
	var field string
	if _, found := moesifOption[fieldName]; found {
		field = moesifOption[fieldName].(func(*http.Request, *http.Response) string)(request, response)
	}
	return field
}

func HeaderToMap(header http.Header) map[string]interface{} {
	headerMap := make(map[string]interface{})
	for name, values := range header {
		headerMap[name] = values
	}
	return headerMap
}

func maskHeaders(headers map[string]interface{}, fieldName string) map[string]interface{} {
	var maskFields []string
	if _, found := moesifOption[fieldName]; found {
		maskFields = moesifOption[fieldName].(func() []string)()
		headers = maskData(headers, maskFields)
	}
	return headers
}

func maskData(data map[string]interface{}, maskBody []string) map[string]interface{} {
	for key, val := range data {
		switch val.(type) {
		case map[string]interface{}:
			if contains(maskBody, key) {
				data[key] = "*****"
			} else {
				maskData(val.(map[string]interface{}), maskBody)
			}
		default:
			if contains(maskBody, key) {
				data[key] = "*****"
			}
		}
	}
	return data
}

func parseBody(readReqBody []byte, fieldName string) (interface{}, string) {
	var body interface{}
	bodyEncoding := "json"
	if jsonMarshalErr := json.Unmarshal(readReqBody, &body); jsonMarshalErr != nil {
		if debug {
			log.Printf("About to parse body as base64 ")
		}
		body = b64.StdEncoding.EncodeToString(readReqBody)
		bodyEncoding = "base64"
		if debug {
			log.Printf("Parsed body as base64 - %s", body)
		}
	} else {
		// Mask Json data
		var maskFields []string
		if _, found := moesifOption[fieldName]; found {
			maskFields = moesifOption[fieldName].(func() []string)()
			body = maskData(body.(map[string]interface{}), maskFields)
		}
	}
	return body, bodyEncoding
}
