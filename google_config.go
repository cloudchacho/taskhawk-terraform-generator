package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
)

// GoogleSchedulerJob struct represents a Taskhawk schedule for periodic jobs
type GoogleSchedulerJob struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	FormatVersion string                 `json:"format_version,omitempty"`
	Headers       map[string]string      `json:"headers,omitempty"`
	Timezone      string                 `json:"timezone,omitempty"`
	Task          string                 `json:"task"`
	Args          []interface{}          `json:"args,omitempty"`
	Kwargs        map[string]interface{} `json:"kwargs,omitempty"`
	Schedule      string                 `json:"schedule"`
}

// GoogleApp struct represents a Taskhawk consumer app
type GoogleApp struct {
	Queue                      string               `json:"queue"`
	ServiceAccounts            []string             `json:"service_accounts"`
	Labels                     map[string]string    `json:"labels"`
	HighMessageCountThresholds map[string]int       `json:"high_message_count_thresholds,omitempty"`
	SchedulerJobs              []GoogleSchedulerJob `json:"scheduler_jobs,omitempty"`
}

// GoogleConfig struct represents the Taskhawk configuration for Google Cloud
type GoogleConfig struct {
	Apps []*GoogleApp `json:"apps,omitempty"`
}

// newGoogleConfig returns a new config read from a file
func newGoogleConfig(filename string) (*GoogleConfig, error) {
	configContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read input file: %q", err)
	}
	config := GoogleConfig{}
	err = json.Unmarshal(configContents, &config)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file as JSON object: %q", err)
	}

	err = config.validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

var googleSubscriptionNameRegex = regexp.MustCompile(`^[a-z0-9-]+$`)
var labelKeyRegex = regexp.MustCompile("^[a-z][a-z0-9-_]*$")
var labelValueRegex = regexp.MustCompile("^[a-z0-9-_]*$")

// Validates that consumer queues are valid format
func (c *GoogleConfig) validateApps() error {
	for _, consumer := range c.Apps {
		if !googleSubscriptionNameRegex.MatchString(consumer.Queue) {
			return fmt.Errorf(
				"invalid subscription name: |%s|, must match regex: %s", consumer.Queue, googleSubscriptionNameRegex)
		}

		for k, v := range consumer.Labels {
			if !labelKeyRegex.MatchString(k) {
				return fmt.Errorf("invalid label key: |%s|, must match regex: %s", k, labelKeyRegex)
			}
			if !labelValueRegex.MatchString(v) {
				return fmt.Errorf("invalid label value: |%s|, must match regex: %s", v, labelValueRegex)
			}
		}
	}
	return nil
}

// validate verifies that the input configuration is sane
func (c *GoogleConfig) validate() error {
	if err := c.validateApps(); err != nil {
		return err
	}

	return nil
}
