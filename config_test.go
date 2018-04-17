package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSchemaFail(t *testing.T) {
	schema := []byte(`
{
  "queue_apps": [
    {
      "queue": "dev-myapp",
      "tags": {
        "App": "myapp",
        "Env": "dev"
      }
    }
  ],
  "lambda_apps": "not-a-list"
}
`)
	assert.EqualError(t, json.Unmarshal(schema, &Config{}),
		"json: cannot unmarshal string into Go struct field Config.lambda_apps of type []*main.LambdaApp")
}

func TestValidateQueue(t *testing.T) {
	invalidQueues := []string{
		"lower",
		"UNDER_SCORE",
		"PUNCTUATION!",
	}

	config := Config{}
	for _, queue := range invalidQueues {
		config.QueueApps = []*QueueApp{{Queue: queue}}
		assert.EqualError(
			t,
			config.validate(),
			fmt.Sprintf("invalid queue name '%s' didn't match: ^[A-Z0-9-]+$", queue),
			queue,
		)
	}
}

func TestValidateNoApps(t *testing.T) {
	config := Config{}
	assert.EqualError(
		t,
		config.validate(),
		"at least one of [QueueApps, LambdaApps] must be non-empty",
	)
}

func TestValidJSON(t *testing.T) {
	var validConfig = []byte(`{
  "queue_apps": [
    {
      "queue": "DEV-MYAPP",
      "tags": {
        "App": "myapp",
        "Env": "dev"
      },
      "schedule": [
        {
          "name": "nightly-job",
          "description": "night job for my app",
          "format_version": "v1.0",
          "headers": {
            "request_id": "1bf0d8a0-5f82-365a-b793-33cac9c9b01d"
          },
          "task": "tasks.send_email",
          "args": [
            "Hello!", 
            "email@automatic.com"
          ],
          "kwargs": {
            "from_email": "spam@example.com"
          }
        }
      ]
    }
  ],
  "lambda_apps": [
    {
      "function_arn": "arn:aws:lambda:us-west-2:12345:function:myFunction:deployed",
      "name": "myapp",
      "schedule": [
        {
          "name": "nightly-job",
          "task": "tasks.cleanup_task"
        }
      ]
    }
  ]
}`)

	var validConfigObj = Config{
		QueueApps: []*QueueApp{
			{
				Queue: "DEV-MYAPP",
				Tags: map[string]string{
					"App": "myapp",
					"Env": "dev",
				},
				Schedule: []ScheduleItem{
					{
						Name:          "nightly-job",
						Description:   "night job for my app",
						FormatVersion: "v1.0",
						Headers: map[string]string{
							"request_id": "1bf0d8a0-5f82-365a-b793-33cac9c9b01d",
						},
						Task: "tasks.send_email",
						Args: []interface{}{
							"Hello!",
							"email@automatic.com",
						},
						Kwargs: map[string]interface{}{
							"from_email": "spam@example.com",
						},
					},
				},
			},
		},
		LambdaApps: []*LambdaApp{
			{
				FunctionARN: "arn:aws:lambda:us-west-2:12345:function:myFunction:deployed",
				Name:        "myapp",
				Schedule: []ScheduleItem{
					{
						Name: "nightly-job",
						Task: "tasks.cleanup_task",
					},
				},
			},
		},
	}

	config := Config{}
	assert.NoError(t, json.Unmarshal(validConfig, &config))
	assert.Equal(t, validConfigObj, config)
}

func TestValidJSONNoLambda(t *testing.T) {
	var validConfig = []byte(`{
  "queue_apps": [
    {
      "queue": "DEV-MYAPP",
      "tags": {
        "App": "myapp",
        "Env": "dev"
      }
    }
  ]
}`)

	var validConfigObj = Config{
		QueueApps: []*QueueApp{
			{
				Queue: "DEV-MYAPP",
				Tags: map[string]string{
					"App": "myapp",
					"Env": "dev",
				},
			},
		},
	}

	config := Config{}
	assert.NoError(t, json.Unmarshal(validConfig, &config))
	assert.Equal(t, validConfigObj, config)
}

func TestValidNoQueueApps(t *testing.T) {
	var validConfig = []byte(`{
  "lambda_apps": [
    {
      "function_arn": "arn:aws:lambda:us-west-2:12345:function:myFunction:deployed",
      "name": "myapp"
    }
  ]
}`)

	var validConfigObj = Config{
		LambdaApps: []*LambdaApp{
			{FunctionARN: "arn:aws:lambda:us-west-2:12345:function:myFunction:deployed", Name: "myapp"},
		},
	}

	config := Config{}
	assert.NoError(t, json.Unmarshal(validConfig, &config))
	assert.Equal(t, validConfigObj, config)
}

func TestLambdaSubscription_Init(t *testing.T) {
	ls := LambdaApp{
		FunctionARN: "arn:aws:lambda:us-west-2:12345:function:myFunction:deployed",
	}
	assert.NoError(t, ls.init())
	assert.Equal(t, "myFunction", ls.FunctionName)
	assert.Equal(t, "deployed", ls.FunctionQualifier)
}

func TestLambdaSubscription_Init_Fail(t *testing.T) {
	ls := LambdaApp{
		FunctionARN: "arn:aws:lambda:us-west-2:12345:foo:myFunction:deployed",
	}
	assert.Error(t, ls.init(), "unable to parse function ARN")
}

func TestLambdaSubscription_Init_NoQualifier(t *testing.T) {
	ls := LambdaApp{
		FunctionARN: "arn:aws:lambda:us-west-2:12345:function:myFunction",
	}
	assert.NoError(t, ls.init())
	assert.Equal(t, "myFunction", ls.FunctionName)
	assert.Equal(t, "", ls.FunctionQualifier)
}

func TestLambdaSubscription_Init_Interpolated(t *testing.T) {
	ls := LambdaApp{
		FunctionARN:  "${aws_lambda_function.myFunction.arn}",
		FunctionName: "myFunction",
	}
	assert.NoError(t, ls.init())
	assert.Equal(t, "myFunction", ls.FunctionName)
	assert.Equal(t, "", ls.FunctionQualifier)
}

func TestLambdaSubscription_Init_InterpolatedFail(t *testing.T) {
	ls := LambdaApp{
		FunctionARN: "${aws_lambda_function.myFunction.arn}",
	}
	assert.Error(t, ls.init(), "unable to parse function ARN since it's an interpolated value")
}
