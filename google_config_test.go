package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateName(t *testing.T) {
	invalidQueues := []string{
		"UPPER",
		"UNDER_SCORE",
		"PUNCTUATION!",
	}

	config := GoogleConfig{}
	for _, queue := range invalidQueues {
		config.Apps = []*GoogleApp{{Queue: queue}}
		assert.EqualError(
			t,
			config.validate(),
			fmt.Sprintf("invalid subscription name: |%s|, must match regex: %s", queue, googleSubscriptionNameRegex),
			"Didn't fail validation for '%s'",
			queue,
		)
	}
}

func TestValidateSubscriptionLabel(t *testing.T) {
	config := GoogleConfig{
		Apps: []*GoogleApp{
			{Queue: "myapp", Labels: map[string]string{"UPPER": ""}},
		},
	}
	assert.EqualError(
		t,
		config.validate(),
		"invalid label key: |UPPER|, must match regex: ^[a-z][a-z0-9-_]*$",
	)
}
