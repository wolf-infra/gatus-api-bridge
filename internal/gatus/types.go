package gatus

// Endpoint matches the structure Gatus expects in config.yaml
type Endpoint struct {
	Name       string            `yaml:"name" json:"name"`
	Group      string            `yaml:"group" json:"group"`
	URL        string            `yaml:"url" json:"url"`
	Method     string            `yaml:"method,omitempty" json:"method,omitempty"`
	Interval   string            `yaml:"interval,omitempty" json:"interval,omitempty"`
	Conditions []string          `yaml:"conditions,omitempty" json:"conditions,omitempty"`
	Headers    map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Alerts     []Alert           `yaml:"alerts,omitempty" json:"alerts,omitempty"`
}

// Alert represents an alert configuration for a specific endpoint
type Alert struct {
	Type             string `yaml:"type" json:"type"`
	FailureThreshold int    `yaml:"failure-threshold,omitempty" json:"failure-threshold,omitempty"`
	SuccessThreshold int    `yaml:"success-threshold,omitempty" json:"success-threshold,omitempty"`
	Description      string `yaml:"description,omitempty" json:"description,omitempty"`
	SendOnResolved   bool   `yaml:"send-on-resolved,omitempty" json:"send-on-resolved,omitempty"`
}

// Config represents the root of the Gatus config file
type Config struct {
	Endpoints []Endpoint `yaml:"endpoints"`
}
