package moesifmiddleware

import (
	"net/http"

	"github.com/moesif/moesifapi-go"
)

type RulesManager interface {
	// Update rules across all stores
	UpdateAllRules(rules []moesifapi.GovernanceRule)

	// Given a mapping of an entity_id, entity type to rule_ids, update the mappings
	UpdateEntityMappings(entityId, entityType string, ruleIds []string)

	// Fetch all rules relevant to an HTTP request from all the stores
	GetAllApplicableRules(req *http.Request) []moesifapi.GovernanceRule
}

type DefaultRulesManager struct {
	stores []RulesStore
}

func NewDefaultRulesManager() *DefaultRulesManager {
	return &DefaultRulesManager{
		stores: []RulesStore{
			RegexRulesStore{},
			EntityRulesStore{Type: "user_id"},
			EntityRulesStore{Type: "company_id"},
		},
	}
}

func (manager *DefaultRulesManager) UpdateAllRules(rules []moesifapi.GovernanceRule) {
	for _, store := range manager.stores {
		store.UpdateRules(rules)
	}
}

func (manager *DefaultRulesManager) UpdateEntityMappings(entityId, entityType string, ruleIds []string) {
	// ... update entity mappings ...
}

func (manager *DefaultRulesManager) GetAllApplicableRules(req *http.Request) []moesifapi.GovernanceRule {
	var allRules []moesifapi.GovernanceRule
	for _, store := range manager.stores {
		allRules = append(allRules, store.GetApplicableRules(req)...)
	}
	return allRules
}

type RulesStore interface {
	// Update the rules in the store
	UpdateRules(rules []moesifapi.GovernanceRule)

	// Get the rules that apply to a given HTTP request
	GetApplicableRules(req *http.Request) []moesifapi.GovernanceRule
}

type RegexRulesStore struct {
	rules []moesifapi.GovernanceRule
}

func (store *RegexRulesStore) UpdateRules(rules []moesifapi.GovernanceRule) {
	store.rules = []moesifapi.GovernanceRule{}
	for _, rule := range rules {
		if rule.Type == "regex" {
			store.rules = append(store.rules, rule)
		}
	}
}

func (store *RegexRulesStore) GetApplicableRules(req *http.Request) []moesifapi.GovernanceRule {
	return store.rules
}

type EntityRulesStore struct {
	Type        string
	RuleCohorts map[string]map[string]bool          // rule_id to entity_id to membership
	Rules       map[string]moesifapi.GovernanceRule // rule_id to rule
}

func NewEntityRulesStore(entityType string) *EntityRulesStore {
	return &EntityRulesStore{
		Type:        entityType,
		RuleCohorts: make(map[string]map[string]bool),
		Rules:       make(map[string]moesifapi.GovernanceRule),
	}
}

func (store *EntityRulesStore) UpdateRules(rules []moesifapi.GovernanceRule) {
	store.Rules = make(map[string]moesifapi.GovernanceRule)
	for _, rule := range rules {
		if rule.Type == store.Type {
			store.Rules[rule.ID] = rule
		}
	}
}

func (store *EntityRulesStore) IsApplicableToEntity(rule moesifapi.GovernanceRule, entity_id string) bool {
	isCohortMember := store.RuleCohorts[rule.ID][entity_id]
	return rule.ApplyTo == "matching" && isCohortMember || rule.ApplyTo == "not_matching" && !isCohortMember
}

func (store *EntityRulesStore) GetApplicableRules(req *http.Request) (rules []moesifapi.GovernanceRule) {
	for _, rule := range store.Rules {
		entity_id := "" // TODO: move entity_id extraction here
		if store.IsApplicableToEntity(rule, entity_id) {
			rules = append(rules, rule.Rule)
		}
	}
	return
}
