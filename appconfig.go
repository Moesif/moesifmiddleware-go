package moesifmiddleware

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

type AppConfig struct {
	Mu      sync.RWMutex
	Updates chan string
	eTags   [2]string
	config  AppConfigResponse
}

func NewAppConfig() AppConfig {
	return AppConfig{
		Updates: make(chan string, 1),
		config:  NewAppConfigResponse(),
	}
}

func (c *AppConfig) Read() AppConfigResponse {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.config
}

func (c *AppConfig) Write(config AppConfigResponse) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.config = config
	c.eTags[1] = c.eTags[0]
	c.eTags[0] = config.eTag
}

func (c *AppConfig) Go() {
	go c.UpdateLoop()
	c.Notify("go")
}

func (c *AppConfig) Notify(eTag string) {
	c.Mu.RLock()
	e := c.eTags
	c.Mu.RUnlock()
	if eTag == "" || eTag == e[0] || eTag == e[1] {
		return
	}
	select {
	case c.Updates <- eTag:
	default:
	}
}

func (c *AppConfig) UpdateLoop() {
	for {
		eTag, more := <-c.Updates
		if !more {
			return
		}
		config, err := getAppConfig()
		if err != nil {
			log.Printf("Failed to get config: %v", err)
			continue
		}
		log.Printf("AppConfig.Notify ETag=%s got /config response ETag=%s", eTag, config.eTag)
		c.Write(config)
	}
}

func (c *AppConfig) GetEntityValues(userId, companyId string) (values []EntityRuleValues) {
	log.Printf("userId=%s companyId=%s", userId, companyId)
	config := c.Read()
	// look up and copy company rules to check
	values = append(values, config.CompanyRules[companyId]...)
	// Lookup and copy user rules to check
	values = append(values, config.UserRules[userId]...)
	return
}

type AppConfigResponse struct {
	OrgID                    string                        `json:"org_id"`
	AppID                    string                        `json:"app_id"`
	SampleRate               int                           `json:"sample_rate"`
	BlockBotTraffic          bool                          `json:"block_bot_traffic"`
	UserSampleRate           map[string]int                `json:"user_sample_rate"`    // user id to a sample rate [0, 100]
	CompanySampleRate        map[string]int                `json:"company_sample_rate"` // company id to a sample rate [0, 100]
	UserRules                map[string][]EntityRuleValues `json:"user_rules"`          // user id to a rule id and template values
	CompanyRules             map[string][]EntityRuleValues `json:"company_rules"`       // company id to a rule id and template values
	IPAddressesBlockedByName map[string]string             `json:"ip_addresses_blocked_by_name"`
	RegexConfig              []RegexRule                   `json:"regex_config"`
	BillingConfigJsons       map[string]string             `json:"billing_config_jsons"`
	eTag                     string
}

func NewAppConfigResponse() AppConfigResponse {
	return AppConfigResponse{
		SampleRate: 100,
	}
}

// EntityRule is a user rule or company rule
type EntityRuleValues struct {
	Rule   string            `json:"rules"`
	Values map[string]string `json:"values"`
}

// Regex Rule
type RegexRule struct {
	Conditions []RegexCondition `json:"conditions"`
	SampleRate int              `json:"sample_rate"`
}

// RegexCondition
type RegexCondition struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

func getAppConfig() (config AppConfigResponse, err error) {
	config = NewAppConfigResponse()
	r, err := apiClient.GetAppConfig()
	if err != nil {
		log.Printf("Application configuration request error: %v", err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Application configuration response body read error: %v", err)
		return
	}
	err = json.Unmarshal(body, &config)
	if err != nil {
		log.Printf("Application configuration response body malformed: %v", err)
		return
	}
	config.eTag = r.Header.Get("X-Moesif-Config-Etag")
	return
}

func getSamplingPercentage(userId string, companyId string) int {
	c := appConfig.Read()
	if userId != "" {
		if userRate, ok := c.UserSampleRate[userId]; ok {
			return userRate
		}
	}

	if companyId != "" {
		if companyRate, ok := c.CompanySampleRate[companyId]; ok {
			return companyRate
		}
	}

	return c.SampleRate
}
