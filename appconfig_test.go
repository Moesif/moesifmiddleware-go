package moesifmiddleware

import (
	"fmt"
	"testing"
)

const id = ""

func TestGetConfig(t *testing.T) {
	moesifClient(map[string]interface{}{
		"Application_Id": id,
		"Api_Endpoint":   "https://api-dev.moesif.net",
	})
	config, err := getAppConfig()
	if err != nil {
		t.Fail()
	}
	fmt.Printf("%#v\n", config)
}

func TestGetRules(t *testing.T) {
	moesifClient(map[string]interface{}{
		"Application_Id": id,
		"Api_Endpoint":   "https://api-dev.moesif.net",
	})
	resp, err := apiClient.GetGovernanceRules()
	if err != nil {
		t.Fail()
	}
	fmt.Printf("%#v\n", resp)
}
