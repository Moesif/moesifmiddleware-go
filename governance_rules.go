package moesifmiddleware

import (
	"log"
	"net/http"
	"regexp"
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
	EntityRules map[string]moesifapi.GovernanceRule
	Regex       []moesifapi.GovernanceRule
}

func NewGovernanceRules() GovernanceRules {
	return GovernanceRules{
		Updates: make(chan string),
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
}

func (g *GovernanceRules) Go() {
	go g.UpdateLoop()
	g.Notify("go")
}

func (g *GovernanceRules) Notify(eTag string) {
	if eTag == "" {
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
		if eTag == g.eTags[0] || eTag == g.eTags[1] {
			continue
		}
		response, err := apiClient.GetGovernanceRules()
		if err != nil {
			continue
		}
		config := NewGovernanceRulesConfig()
		for _, r := range response.Rules {
			switch r.Type {
			case "user", "company":
				config.EntityRules[r.ID] = r
			case "regexp":
				config.Regex = append(config.Regex, r)
			}

		}
		g.Mu.Lock()
		g.Write(config)
		g.Mu.Unlock()
		g.eTags[1] = g.eTags[0]
		g.eTags[0] = response.ETag
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

func (g *GovernanceRules) Get(request *http.Request, entityValues []EntityRuleValues) (rules []RuleTemplate) {
	config := g.Read()
	// in a list of rules with overrides, the last override value is what will be used in the response
	// create a slice of rules to check in reverse priority order
	// regex rule, company rule, user rule order, i.e. user rule overrides take priority over company, etc.
	regexToCheck := make([]RuleTemplate, len(config.Regex))
	// copy all regex rules first
	for i, r := range config.Regex {
		regexToCheck[i] = RuleTemplate{Rule: r}
	}
	//copy all user and company entity rules
	for _, ev := range entityValues {
		if rule, ok := config.EntityRules[ev.Rule]; ok {
			regexToCheck = append(regexToCheck, RuleTemplate{rule, ev.Values})
		}
	}
	// if a rule from above has regex conditions, the rule is used if matching; otherwise, it's used
	for _, r := range regexToCheck {
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

func RequestPathLookup(req *http.Request, path string) string {
	switch path {
	case "request.ip_address":
		return req.RemoteAddr
	case "request.route":
		return req.URL.Path
	case "request.verb":
		return req.Method
	default:
		return ""
	}
}

func CheckRegex(rule moesifapi.GovernanceRule, req *http.Request) bool {
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
