package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// ScheduleItem struct represents a Taskhawk schedule for periodic jobs
type ScheduleItem struct {
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	FormatVersion      string                 `json:"format_version,omitempty"`
	Headers            map[string]string      `json:"headers,omitempty"`
	Task               string                 `json:"task"`
	Args               []interface{}          `json:"args,omitempty"`
	Kwargs             map[string]interface{} `json:"kwargs,omitempty"`
	ScheduleExpression string                 `json:"schedule_expression"`
}

func isValidPriority(value string) bool {
	return value == "default" || value == "high" || value == "bulk" || value == "low"
}

// QueueApp struct represents a Taskhawk consumer app
type QueueApp struct {
	Queue                      string            `json:"queue"`
	Tags                       map[string]string `json:"tags"`
	Schedule                   []ScheduleItem    `json:"schedule,omitempty"`
	HighMessageCountThresholds map[string]int    `json:"high_message_count_thresholds,omitempty"`
}

// LambdaApp struct represents a Taskhawk subscription for a lambda app
type LambdaApp struct {
	FunctionARN                string         `json:"function_arn"`
	FunctionName               string         `json:"function_name,omitempty"`
	FunctionQualifier          string         `json:"function_qualifier,omitempty"`
	Name                       string         `json:"name"`
	Schedule                   []ScheduleItem `json:"schedule,omitempty"`
	HighMessageCountThresholds map[string]int `json:"high_message_count_thresholds,omitempty"`
}

var lambdaARNRegexp = regexp.MustCompile(`^arn:aws:lambda:([^:]+):([^:]+):function:([^:]+)(:([^:]+))?$`)

// init initializes the data structure with function name and qualifier if required
func (ls *LambdaApp) init() error {
	if ls.FunctionName != "" {
		return nil
	}

	if strings.Contains(ls.FunctionARN, "${") {
		return fmt.Errorf("unable to parse function ARN since it's an interpolated value")
	}

	matches := lambdaARNRegexp.FindStringSubmatch(ls.FunctionARN)
	if len(matches) > 0 {
		if ls.FunctionName == "" {
			ls.FunctionName = matches[3]
		}
		if ls.FunctionQualifier == "" && len(matches) >= 6 {
			ls.FunctionQualifier = matches[5]
		}
	}
	if ls.FunctionName == "" {
		return fmt.Errorf("unable to parse function ARN")
	}
	return nil
}

// Config struct represents the Taskhawk configuration
type Config struct {
	QueueApps  []*QueueApp  `json:"queue_apps,omitempty"`
	LambdaApps []*LambdaApp `json:"lambda_apps,omitempty"`
}

// NewConfig returns a new config read from a file
func NewConfig(filename string) (*Config, error) {
	configContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read input file: %q", err)
	}
	config := Config{}
	err = json.Unmarshal(configContents, &config)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read config file as JSON object")
	}

	err = config.validate()
	if err != nil {
		return nil, err
	}

	count := len(config.LambdaApps)
	for i, ls := range config.LambdaApps {
		err := ls.init()
		if err != nil {
			return nil, errors.Wrapf(err, "init failure on app %d/%d", i, count)
		}
	}
	return &config, nil
}

var lambdaAppNameRegex = regexp.MustCompile(`^[a-z0-9-]+$`)
var queueRegex = regexp.MustCompile(`^[A-Z0-9-]+$`)

// Validates that consumer queues are valid format
func (c *Config) validateConsumers() error {
	for _, consumer := range c.QueueApps {
		if !queueRegex.MatchString(consumer.Queue) {
			return fmt.Errorf("invalid queue name '%s' didn't match: %s", consumer.Queue, queueRegex)
		}

		if consumer.HighMessageCountThresholds != nil {
			for priority, threshold := range consumer.HighMessageCountThresholds {
				if !isValidPriority(priority) {
					return errors.Errorf("invalid priority: '%s'", priority)
				}
				if threshold < 0 {
					return errors.Errorf("invalid threshold: '%d'", threshold)
				}
			}
		}
	}
	return nil
}

// Validates that lambda subscriptions refer to valid topics
func (c *Config) validateLambdaApps() error {
	for _, app := range c.LambdaApps {
		if !lambdaAppNameRegex.MatchString(app.Name) {
			return fmt.Errorf("invalid lambda app name '%s' didn't match: %s", app.Name, lambdaAppNameRegex)
		}
	}
	return nil
}

// validate verifies that the input configuration is sane
func (c *Config) validate() error {
	if err := c.validateConsumers(); err != nil {
		return err
	}

	if err := c.validateLambdaApps(); err != nil {
		return err
	}

	if len(c.LambdaApps) == 0 && len(c.QueueApps) == 0 {
		return fmt.Errorf("at least one of [QueueApps, LambdaApps] must be non-empty")
	}

	return nil
}
