package analyzer

import (
    "sync"
)

type Registry struct {
    rules map[string]Rule
    mu    sync.RWMutex
}

var defaultRegistry = NewRegistry()

func NewRegistry() *Registry {
    return &Registry{
        rules: make(map[string]Rule),
    }
}

func (r *Registry) Register(rule Rule) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.rules[rule.Name()] = rule
}

func (r *Registry) GetRule(name string) (Rule, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    rule, ok := r.rules[name]
    return rule, ok
}

func (r *Registry) GetAllRules() []Rule {
    r.mu.RLock()
    defer r.mu.RUnlock()
    rules := make([]Rule, 0, len(r.rules))
    for _, rule := range r.rules {
        rules = append(rules, rule)
    }
    return rules
}

func Register(rule Rule) {
    defaultRegistry.Register(rule)
}

func GetAllRules() []Rule {
    return defaultRegistry.GetAllRules()
}
