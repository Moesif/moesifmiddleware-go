package moesifmiddleware

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func parseConfiguration(response *http.Response) (int, string, time.Time) {

	// Get X-Moesif-Config-Etag header from response
	if configETag, ok := response.Header["X-Moesif-Config-Etag"]; ok {
		eTag = configETag[0]
	}

	// Read the response body
	readRespBody, err := ioutil.ReadAll(response.Body)
	if err == nil {
		// Parse the response Body
		if jsonRespParseErr := json.Unmarshal(readRespBody, &appConfig); jsonRespParseErr == nil {
			// Fetch sample rate from appConfig
			if getSampleRate, found := appConfig["sample_rate"]; found {
				if rate, ok := getSampleRate.(float64); ok {
					samplingPercentage = int(rate)
				}
			}
			// Fetch User Sample rate from appConfig
			if userSampleRate, ok := appConfig["user_sample_rate"]; ok {
				if userRates, ok := userSampleRate.(map[string]interface{}); ok {
					userSampleRateMap = userRates
				}
			}
			// Fetch Company Sample rate from appConfig
			if companySampleRate, ok := appConfig["company_sample_rate"]; ok {
				if companyRates, ok := companySampleRate.(map[string]interface{}); ok {
					companySampleRateMap = companyRates
				}
			}
		}
	}

	return samplingPercentage, eTag, time.Now().UTC()
}

func getSamplingPercentage(userId string, companyId string) int {

	if userId != "" {
		if userRate, ok := userSampleRateMap[userId].(float64); ok {
			return int(userRate)
		}
	}

	if companyId != "" {
		if companyRate, ok := companySampleRateMap[companyId].(float64); ok {
			return int(companyRate)
		}
	}

	if getSampleRate, found := appConfig["sample_rate"]; found {
		if rate, ok := getSampleRate.(float64); ok {
			return int(rate)
		}
	}

	return 100
}
