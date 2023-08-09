package moesifmiddleware

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/moesif/moesifapi-go"
)

type GovernanceRules struct {
	Mu      sync.RWMutex
	Updates chan string
	eTags   [2]string
	config  GovernanceRulesConfig
}

type GovernanceRulesConfig struct {
	EntityRules  map[string]moesifapi.GovernanceRule
	UserRules    []moesifapi.GovernanceRule
	CompanyRules []moesifapi.GovernanceRule
	Regex        []moesifapi.GovernanceRule
	eTag         string
}

func NewGovernanceRules() GovernanceRules {
	return GovernanceRules{
		Updates: make(chan string, 1),
		config:  NewGovernanceRulesConfig(),
	}
}

func (g *GovernanceRules) Read() GovernanceRulesConfig {
	g.Mu.RLock()
	defer g.Mu.RUnlock()
	return g.config
}

func (g *GovernanceRules) Write(config GovernanceRulesConfig) {
	g.Mu.Lock()
	defer g.Mu.Unlock()
	g.config = config
	g.eTags[1] = g.eTags[0]
	g.eTags[0] = config.eTag
}

func (g *GovernanceRules) Go() {
	go g.UpdateLoop()
	g.Notify("go")
}

func (g *GovernanceRules) Notify(eTag string) {
	g.Mu.RLock()
	e := g.eTags
	g.Mu.RUnlock()
	if eTag == "" || eTag == e[0] || eTag == e[1] {
		return
	}
	select {
	case g.Updates <- eTag:
	default:
	}
}

func (g *GovernanceRules) UpdateLoop() {
	for {
		eTag, more := <-g.Updates
		if !more {
			return
		}
		response, err := apiClient.GetGovernanceRules()
		if err != nil {
			log.Printf("Failed to get governance rules: %v", err)
			continue
		}
		config := NewGovernanceRulesConfig()
		config.eTag = response.ETag
		for _, r := range response.Rules {
			switch r.Type {
			case "user":
				config.UserRules = append(config.UserRules, r)
				config.EntityRules[r.ID] = r
			case "company":
				config.CompanyRules = append(config.CompanyRules, r)
				config.EntityRules[r.ID] = r
			case "regex":
				config.Regex = append(config.Regex, r)
			}
		}
		log.Printf("GovernanceRules.Notify ETag=%s got /rules response ETag=%s", eTag, config.eTag)
		g.Write(config)
	}
}

func NewGovernanceRulesConfig() (g GovernanceRulesConfig) {
	g.EntityRules = make(map[string]moesifapi.GovernanceRule)
	return
}

type RuleTemplate struct {
	Rule   moesifapi.GovernanceRule
	Values map[string]string
}

func (r RuleTemplate) TemplateOverride() (t TemplatedOverrideValues) {
	t.Block = r.Rule.Block
	t.Status = r.Rule.ResponseOverrides.Status
	t.Headers = make(map[string]string)
	for k, v := range r.Rule.ResponseOverrides.Headers {
		t.Headers[k] = moesifapi.Template(v, r.Values)
	}
	t.Body = []byte(moesifapi.Template(string(r.Rule.ResponseOverrides.Body), r.Values))
	return
}

type TemplatedOverrideValues struct {
	Block   bool
	Headers map[string]string
	Status  int
	Body    []byte
}

func GetMatchingRuleTemplates(config GovernanceRulesConfig, values []EntityRuleValues) (inCohorts map[string]bool, templates []RuleTemplate) {
	inCohorts = make(map[string]bool)
	for _, ev := range values {
		// check all the rules for which this entity is in a cohort
		inCohorts[ev.Rule] = true
		// if the rule is a matching rule, add it to the list of templates
		if rule, ok := config.EntityRules[ev.Rule]; ok && rule.ApplyTo == "matching" {
			templates = append(templates, RuleTemplate{rule, ev.Values})
		}
	}
	return
}

func GetNotMatchingRuleTemplates(rules []moesifapi.GovernanceRule, entityId string, isInCohort map[string]bool) (templates []RuleTemplate) {
	for _, rule := range rules {
		if rule.ApplyTo == "not_matching" && !isInCohort[rule.ID] || rule.ApplyUnidentified && entityId == "" {
			templates = append(templates, RuleTemplate{Rule: rule})
		}
	}
	return
}

func (g *GovernanceRules) Get(request *http.Request, userValues, companyValues []EntityRuleValues, userId, companyId string) (rules []RuleTemplate) {
	config := g.Read()
	// in a list of rules with overrides, the last override value is what will be used in the response
	// create a slice of rules to check in priority order
	// user rule order, company rule, regex rule, i.e. user rule overrides take priority over company, etc.
	regexToCheck := []RuleTemplate{}

	// if a user_id is matching a cohort in any rule, it will have an EntityRuleValues entry in userValues
	// collecting the rule ids that match the company_id in isInCohort allows us to efficiently
	// and simply check if the current entity is in any cohort for a given rule to apply non_matching rules
	isInCohort, matching := GetMatchingRuleTemplates(config, userValues)
	regexToCheck = append(regexToCheck, matching...)
	regexToCheck = append(regexToCheck, GetNotMatchingRuleTemplates(config.UserRules, userId, isInCohort)...)

	isInCohort, matching = GetMatchingRuleTemplates(config, companyValues)
	regexToCheck = append(regexToCheck, matching...)
	regexToCheck = append(regexToCheck, GetNotMatchingRuleTemplates(config.CompanyRules, companyId, isInCohort)...)

	for _, r := range config.Regex {
		regexToCheck = append(regexToCheck, RuleTemplate{Rule: r})
	}

	// if a rule from above has regex conditions, the rule is used if matching; otherwise, it's used
	// the rules have priority in order from highest to lowest.  We apply the list in reverse order so that
	// the highest priority rules are applied last and thus their value is used in the final response
	for i := len(regexToCheck) - 1; i >= 0; i-- {
		r := regexToCheck[i]
		if CheckRegex(r.Rule, request) {
			rules = append(rules, r)
		}
	}
	return
}

type ResponseOverride struct {
	http.ResponseWriter
	Override     TemplatedOverrideValues
	wroteHeaders bool
	wroteBody    bool
}

func NewResponseOverride(response http.ResponseWriter, templates []RuleTemplate) (r ResponseOverride) {
	r.Override.Headers = make(map[string]string)
	r.ResponseWriter = response
	for _, t := range templates {
		o := t.TemplateOverride()
		if o.Block {
			r.Override.Block = true
		}
		if o.Status != 0 {
			r.Override.Status = o.Status
		}
		for k, v := range o.Headers {
			r.Override.Headers[k] = v
		}
		if len(o.Body) > 0 {
			r.Override.Body = o.Body
		}
	}
	return
}

func (r *ResponseOverride) WriteHeader(status int) {
	r.wroteHeaders = true
	h := r.Header()
	for k, v := range r.Override.Headers {
		h.Set(k, v)
	}
	if r.Override.Block {
		status = r.Override.Status
	}
	r.ResponseWriter.WriteHeader(status)
}

func (r *ResponseOverride) Write(body []byte) (int, error) {
	r.wroteBody = true
	if !r.wroteHeaders {
		r.WriteHeader(200)
	}
	if r.Override.Block {
		body = []byte(r.Override.Body)
	}
	return r.ResponseWriter.Write(body)
}

func (r *ResponseOverride) finish() {
	if !r.wroteBody {
		r.Write([]byte{})
	}
}

// bufferRequestBody reads the request body into a buffer and update the request
// object with a reader containing a copy of the contents so that the request
// may be used as normal
func bufferRequestBody(req *http.Request) (body []byte) {
	newBody, b, err := teeBody(req.Body)
	if err != nil {
		log.Printf("Unable to read incoming request body: %v", err)
		return
	}
	req.Body = newBody
	body, _ = ioutil.ReadAll(b)
	return
}

// getJsonKeyString attempts to read a JSON object from input j and
// if object.key exists and has a string value, it is returned
func getJsonKeyString(j []byte, key string) (s string) {
	d := make(map[string]json.RawMessage)
	if json.Unmarshal(j, &d) == nil {
		if m, ok := d[key]; ok {
			json.Unmarshal(m, &s)
		}
	}
	return
}

// RequestPathLookup returns the string in a given regexp matching path from req
func RequestPathLookup(req *http.Request, path string) string {
	switch path {
	case "request.ip_address":
		return req.RemoteAddr
	case "request.route":
		return req.URL.Path
	case "request.verb":
		return req.Method
	}
	const requestBody = "request.body."
	t := req.Header.Get("Content-Type")
	if strings.HasPrefix(path, requestBody) && (t == "application/graphql" || t == "application/json") {
		body := bufferRequestBody(req)
		key := strings.TrimPrefix(path, requestBody)
		if t == "application/graphql" && key == "query" {
			return string(body)
		}
		if t == "application/json" {
			return getJsonKeyString(body, key)
		}
	}
	return ""
}

func CheckRegex(rule moesifapi.GovernanceRule, req *http.Request) bool {
	// if no regex conditions are specified, the rule matches
	if len(rule.RegexConfigOr) == 0 {
		return true
	}
	// a slice of slices of regular expressions is matched
	// the top level slice is logically OR compared, returning true if any inner slices eval true
	// the inner level slices of regular expressions are logically AND compared, only returning true
	// if all expressions in a single inner slice match
	for _, regexAnd := range rule.RegexConfigOr {
		andValue := true
		for _, c := range regexAnd.Conditions {
			s := RequestPathLookup(req, c.Path)
			// c.Value is a regular expression, but if it contains an error, default to false.
			// False here will fail to match the rule which errors on the side of propagating the event
			// rather than a regex error potentially causing a rule to match
			match, err := regexp.MatchString(c.Value, s)
			if err != nil {
				log.Printf(`Governance rule regexp error: org-app=%s-%s rule.id=%s rule.name="%s" path=%s regexp="%s"`, rule.OrgID, rule.AppID, rule.ID, rule.Name, c.Path, c.Value)
			}
			andValue = andValue && match
		}
		if andValue {
			return true
		}
	}
	return false
}
