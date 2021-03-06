package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1" // imports as package "cli"
)

const (
	cloudProviderGoogle = "google"
	cloudProviderAWS    = "aws"
)

const (
	// VERSION represents the version of the generator tool
	VERSION = "v4.2.5"

	// TFAWSQueueModuleVersion represents the version of the AWS taskhawk-queue module
	TFAWSQueueModuleVersion = "1.1.4"

	// TFAWSLambdaModuleVersion represents the version of the AWS taskhawk-lambda module
	TFAWSLambdaModuleVersion = "1.1.1"

	// TFAWSSchedulerModuleVersion represents the version of the AWS taskhawk-scheduler module
	TFAWSSchedulerModuleVersion = "1.0.6"

	// TFGoogleQueueModuleVersion represents the version of the Google taskhawk-queue module
	TFGoogleQueueModuleVersion = "3.2.4"

	tfDoNotEditStamp = `// DO NOT EDIT
// This file has been auto-generated by taskhawk-terraform-generator ` + VERSION
)

const (
	// alertingFlag represents the cli flag that indicates if alerting should be generated
	alertingFlag = "alerting"

	// awsAccountIDFlag represents the cli flag for aws account id
	awsAccountIDFlag = "aws-account-id"

	// awsRegionFlag represents the cli flag for aws region
	awsRegionFlag = "aws-region"

	// cloudProviderFlag represents the cli flag for cloud provider name
	cloudProviderFlag = "cloud"

	// dataflowTmpGCSLocationFlag represents the cli flag for Dataflow temporary GCS location (Google only)
	dataflowTmpGCSLocationFlag = "dataflow-tmp-gcs-location"

	// dataflowPubSubToPubSubTemplateGCSPathFlag represents the cli flag for Dataflow template GCS path
	// for pub sub to pub sub dataflow (Google only)
	dataflowPubSubToPubSubTemplateGCSPathFlag = "dataflow-template-pubsub-to-pubsub-gcs-path"

	// dataflowPubSubToStorageGCSPathFlag represents the cli flag for Dataflow template GCS path
	// for pub sub to Storage dataflow (Google only)
	dataflowPubSubToStorageGCSPathFlag = "dataflow-template-pubsub-to-storage-gcs-path"

	// dataflowAlertNotificationChannelsFlag represents the cli flag for Dataflow freshness alert notification
	// channels (Google only)
	dataflowAlertNotificationChannelsFlag = "dataflow-alert-notification-channels"

	// dlqAlertAlarmActionsFlag represents the cli flag for DLQ alert actions on ALARM
	dlqAlertAlarmActionsFlag = "dlq-alert-alarm-actions"

	// dlqAlertOKActionsFlag represents the cli flag for DLQ alert notification channels (Google only)
	dlqAlertNotificationChannelsFlag = "dlq-alert-notification-channels"

	// dlqAlertOKActionsFlag represents the cli flag for DLQ alert actions on OK
	dlqAlertOKActionsFlag = "dlq-alert-ok-actions"

	// enableFirehoseAllMessages represents the cli flag to enable Google Firehose for all taskhawk messages (Google only
	// for now)
	enableFirehoseAllMessages = "enable-firehose-all-messages"

	// googleDataflowZoneFlag represents the cli flag for Dataflow template GCS zone (Google only)
	googleDataflowZoneFlag = "dataflow-zone"

	// googleDataflowRegionFlag represents the cli flag for Dataflow template GCS region (Google only)
	googleDataflowRegionFlag = "dataflow-region"

	// googleFirehoseDataflowOutputDirectoryFlag represents the cli flag for Firehose Dataflow output directory
	// (Google only)
	googleFirehoseDataflowOutputDirectoryFlag = "google-firehose-dataflow-output-dir"

	// googleProjectAlertingFlag represents the cli flag that indicates the google project for alerting resources
	googleProjectAlertingFlag = "google-project-alerting"

	// iamFlag represents the cli flag for iam generation
	iamFlag = "iam"

	// moduleFlag represents the cli flag for output module name
	moduleFlag = "module"

	// queueAlertAlarmActionsFlag represents the cli flag for DLQ alert actions on ALARM
	queueAlertAlarmActionsFlag = "queue-alert-alarm-actions"

	// queueAlertNotificationChannelsFlag represents the cli flag for queue alert notification channels (Google only)
	queueAlertNotificationChannelsFlag = "queue-alert-notification-channels"

	// queueAlertOKActionsFlag represents the cli flag for DLQ alert actions on OK
	queueAlertOKActionsFlag = "queue-alert-ok-actions"

	// highMessageCountThresholdFlag represents the cli flag for high message count
	highMessageCountThresholdFlag = "high-message-count-threshold"
)

var providerSpecificFlags = map[string][]string{
	cloudProviderAWS: {
		awsAccountIDFlag,
		awsRegionFlag,
		dlqAlertAlarmActionsFlag,
		dlqAlertOKActionsFlag,
		queueAlertAlarmActionsFlag,
		queueAlertOKActionsFlag,
	},
	cloudProviderGoogle: {
		dlqAlertNotificationChannelsFlag,
		dataflowPubSubToPubSubTemplateGCSPathFlag,
		dataflowPubSubToStorageGCSPathFlag,
		dataflowTmpGCSLocationFlag,
		enableFirehoseAllMessages,
		googleDataflowZoneFlag,
		googleDataflowRegionFlag,
		googleFirehoseDataflowOutputDirectoryFlag,
		googleProjectAlertingFlag,
		queueAlertNotificationChannelsFlag,
	},
}

var providerRequiredFlags = map[string][]string{
	cloudProviderAWS: {
		awsAccountIDFlag,
		awsRegionFlag,
	},
	cloudProviderGoogle: {},
}

var providerAlertingFlags = map[string][]string{
	cloudProviderAWS: {
		queueAlertAlarmActionsFlag,
		queueAlertOKActionsFlag,
		dlqAlertAlarmActionsFlag,
		dlqAlertOKActionsFlag,
	},
	cloudProviderGoogle: {
		queueAlertNotificationChannelsFlag,
		dlqAlertNotificationChannelsFlag,
		dataflowAlertNotificationChannelsFlag,
		googleProjectAlertingFlag,
	},
}

var providerAlertingRequiredFlags = map[string][]string{
	cloudProviderAWS: {
		queueAlertAlarmActionsFlag,
		queueAlertOKActionsFlag,
		dlqAlertAlarmActionsFlag,
		dlqAlertOKActionsFlag,
	},
	cloudProviderGoogle: {
		queueAlertNotificationChannelsFlag,
		dlqAlertNotificationChannelsFlag,
		dataflowAlertNotificationChannelsFlag,
	},
}

func validateArgs(c *cli.Context) *cli.ExitError {
	cloudProvider := c.GlobalString(cloudProviderFlag)
	if cloudProvider == "" {
		return cli.NewExitError(fmt.Sprintf("--%s is required", cloudProviderFlag), 1)
	}
	if cloudProvider != cloudProviderAWS && cloudProvider != cloudProviderGoogle {
		return cli.NewExitError(fmt.Sprintf("invalid cloud provider: %s", cloudProvider), 1)
	}

	if c.NArg() != 1 {
		return cli.NewExitError("<config-file> is required", 1)
	}

	// verify provider flags are used correctly
	for provider, flags := range providerSpecificFlags {
		if provider == cloudProvider {
			continue
		}
		for _, flag := range flags {
			if c.IsSet(flag) {
				return cli.NewExitError(
					fmt.Sprintf("flag --%s disallowed for provider: %s", flag, cloudProvider),
					1,
				)
			}
		}
	}

	// verify required flags are provided
	for _, flag := range providerRequiredFlags[cloudProvider] {
		if !c.IsSet(flag) {
			return cli.NewExitError(
				fmt.Sprintf("flag --%s is required for provider: %s", flag, cloudProvider),
				1,
			)
		}
	}

	// verify alerting flags are used correctly
	alertingFlagsOkay := true
	if c.Bool(alertingFlag) {
		for _, f := range providerAlertingRequiredFlags[cloudProvider] {
			if !c.IsSet(f) {
				alertingFlagsOkay = false
				msg := fmt.Sprintf("--%s is required\n", f)
				if _, err := fmt.Fprint(cli.ErrWriter, msg); err != nil {
					return cli.NewExitError(msg, 1)
				}
			}
		}
		if !alertingFlagsOkay {
			return cli.NewExitError("missing required flags for --alerting", 1)
		}
	} else {
		for _, f := range providerAlertingFlags[cloudProvider] {
			if c.IsSet(f) {
				alertingFlagsOkay = false
				msg := fmt.Sprintf("--%s is disallowed\n", f)
				if _, err := fmt.Fprint(cli.ErrWriter, msg); err != nil {
					return cli.NewExitError(msg, 1)
				}
			}
		}
		if !alertingFlagsOkay {
			return cli.NewExitError("disallowed flags specified with missing --alerting", 1)
		}
	}

	return nil
}

func generateModule(c *cli.Context) error {
	if err := validateArgs(c); err != nil {
		return err
	}

	configFile := c.Args().Get(0)

	config, err := newConfig(c, configFile)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	writer := newConfigWriter(c, config)
	err = writer.writeTerraform()
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "failed to generate terraform module"), 1)
	}

	fmt.Println("Created Terraform Taskhawk module successfully!")
	return nil
}

func generateConfigFileStructure(c *cli.Context) error {
	cloudProvider := c.GlobalString(cloudProviderFlag)
	if cloudProvider == "" {
		return cli.NewExitError(fmt.Sprintf("--%s is required", cloudProviderFlag), 1)
	}
	if cloudProvider != cloudProviderAWS && cloudProvider != cloudProviderGoogle {
		return cli.NewExitError(fmt.Sprintf("invalid cloud provider: %s", cloudProvider), 1)
	}

	var structure interface{}
	if cloudProvider == cloudProviderAWS {
		structure = AWSConfig{
			QueueApps: []*AWSQueueApp{
				{
					Queue: "DEV-MYAPP",
					Tags: map[string]string{
						"App": "myapp",
						"Env": "dev",
					},
					Schedule: []AWSScheduleItem{
						{
							Name:          "nightly-job (unique for each app)",
							Description:   "{optional description}",
							FormatVersion: "{optional format version}",
							Headers: map[string]string{
								"header": "{optional headers}",
							},
							Task: "tasks.send_email",
							Args: []interface{}{"{optional args}"},
							Kwargs: map[string]interface{}{
								"kwarg1": "{optional keyword args}",
							},
							ScheduleExpression: "*/10 * * * *",
						},
					},
					HighMessageCountThresholds: map[string]int{
						"bulk": 100000,
					},
				},
			},
			LambdaApps: []*LambdaApp{
				{
					FunctionARN:  "arn:aws:lambda:us-west-2:12345:function:my_function:deployed",
					FunctionName: "{optional - this value is inferred from FunctionARN if that's not an interpolated value}",
					FunctionQualifier: "{optional - this value is inferred from FunctionARN if that's not an interpolated" +
						" value}",
					Name: "myapp",
					Schedule: []AWSScheduleItem{
						{
							Name:          "nightly-job (unique for each app)",
							Description:   "{optional description}",
							FormatVersion: "{optional format version}",
							Headers: map[string]string{
								"header": "{optional headers}",
							},
							Task: "tasks.send_email",
							Args: []interface{}{"{optional args}"},
							Kwargs: map[string]interface{}{
								"kwarg1": "{optional keyword args}",
							},
							ScheduleExpression: "*/10 * * * *",
						},
					},
				},
			},
		}
	} else if cloudProvider == cloudProviderGoogle {
		structure = GoogleConfig{
			Apps: []*GoogleApp{
				{
					"dev-myapp",
					[]string{
						"myapp@project.iam.gserviceaccount.com",
					},
					map[string]string{
						"epp": "myapp",
						"env": "dev",
					},
					map[string]int{
						"bulk": 100000,
					},
					[]GoogleSchedulerJob{
						{
							Name:          "nightly-job (unique for each app)",
							Description:   "{optional description}",
							Priority:      "high",
							FormatVersion: "{optional format version}",
							Headers: map[string]string{
								"header": "{optional headers}",
							},
							Timezone: "America/Los_Angeles",
							Task:     "tasks.send_email",
							Args:     []interface{}{"{optional args}"},
							Kwargs: map[string]interface{}{
								"kwarg1": "{optional keyword args}",
							},
							Schedule: "*/10 * * * *",
						},
					},
				},
			},
		}
	}
	structureAsJSON, err := json.MarshalIndent(structure, "", "    ")
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	fmt.Println(string(structureAsJSON))
	return nil
}

func runApp(args []string) error {
	cli.VersionFlag = cli.BoolFlag{Name: "version, V"}

	app := cli.NewApp()
	app.Name = "TaskHawk Terraform"
	app.Usage = "Manage Terraform configuration for TaskHawk apps"
	app.Version = VERSION
	app.HelpName = "taskhawk-terraform"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  cloudProviderFlag,
			Usage: "Cloud provider - either aws or google",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "generate",
			Usage:     "Generates Terraform module for TaskHawk apps",
			ArgsUsage: "<config-file>",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  moduleFlag,
					Usage: "Terraform module name to generate",
					Value: "taskhawk",
				},
				cli.BoolFlag{
					Name:  alertingFlag,
					Usage: "Should alerting be generated?",
				},
				cli.BoolFlag{
					Name:  iamFlag,
					Usage: "Should IAM policies be generated? (AWS only)",
				},
				cli.StringSliceFlag{
					Name:  queueAlertAlarmActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in queue when in ALARM (AWS only)",
				},
				cli.StringSliceFlag{
					Name:  queueAlertOKActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in queue when OK (AWS only)",
				},
				cli.StringSliceFlag{
					Name:  dlqAlertAlarmActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in dead-letter queue when in ALARM (AWS only)",
				},
				cli.StringSliceFlag{
					Name:  dlqAlertOKActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in dead-letter queue when OK (AWS only)",
				},
				cli.StringFlag{
					Name:  awsAccountIDFlag,
					Usage: "AWS Account ID (AWS only)",
				},
				cli.StringFlag{
					Name:  awsRegionFlag,
					Usage: "AWS Region (AWS only)",
				},
				cli.IntFlag{
					Name:  highMessageCountThresholdFlag,
					Usage: "High message count threshold",
				},
				cli.StringFlag{
					Name:  dataflowTmpGCSLocationFlag,
					Usage: "Dataflow tmp GCS location (Google only)",
				},
				cli.StringFlag{
					Name:  dataflowPubSubToPubSubTemplateGCSPathFlag,
					Usage: "Dataflow template for pubsub to pubsub GCS location (Google only)",
				},
				cli.StringFlag{
					Name:  dataflowPubSubToStorageGCSPathFlag,
					Usage: "Dataflow template for pubsub to storage GCS location (Google only)",
				},
				cli.StringSliceFlag{
					Name:  dataflowAlertNotificationChannelsFlag,
					Usage: "Dataflow freshness alert notification channels (Google only)",
				},
				cli.BoolFlag{
					Name:  enableFirehoseAllMessages,
					Usage: "Enable Google Firehose for all taskhawk messages (Google only for now)",
				},
				cli.StringFlag{
					Name:  googleFirehoseDataflowOutputDirectoryFlag,
					Usage: "Google Firehose Dataflow output directory. Must end with /. (Google only)",
				},
				cli.StringFlag{
					Name: googleDataflowZoneFlag,
					Usage: "Dataflow zone (Google only) (required if zone isn't set at provider level, or " +
						"isn't supported by Dataflow)",
				},
				cli.StringFlag{
					Name: googleDataflowRegionFlag,
					Usage: "Dataflow region (Google only) (required if region isn't set at provider level, or " +
						"isn't supported by Dataflow)",
				},
				cli.StringSliceFlag{
					Name:  dlqAlertNotificationChannelsFlag,
					Usage: "Stackdriver Notification Channels for high message count in dead-letter queue (Google only)",
				},
				cli.StringSliceFlag{
					Name:  queueAlertNotificationChannelsFlag,
					Usage: "Stackdriver Notification Channels for high message count in queue (Google only)",
				},
				cli.StringFlag{
					Name: googleProjectAlertingFlag,
					Usage: "Google project to use for alerting resources. This may be different from your main" +
						"app environment (Google only)",
				},
			},
			Action: generateModule,
		},
		{
			Name:   "config-file-structure",
			Usage:  "Outputs the structure for config file required for generate command",
			Action: generateConfigFileStructure,
		},
	}

	return app.Run(args)
}

func main() {
	if err := runApp(os.Args); err != nil {
		log.Fatal(err)
	}
}
