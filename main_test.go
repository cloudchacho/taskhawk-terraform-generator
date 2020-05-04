package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func argsForTestNoOptional(cloudProvider string, configFilepath string) []string {
	args := []string{
		"./taskhawk-terraform-generator",
		fmt.Sprintf("--%s", cloudProviderFlag),
		cloudProvider,
		"generate",
		configFilepath,
	}
	if cloudProvider == cloudProviderAWS {
		args = append(
			args,
			fmt.Sprintf(`--%s=12345`, awsAccountIDFlag),
			fmt.Sprintf(`--%s=us-west-2`, awsRegionFlag),
		)
	}
	return args
}

func argsForTest(cloudProvider string, configFilepath string) []string {
	args := argsForTestNoOptional(cloudProvider, configFilepath)
	args = append(
		args,
		"--alerting",
		fmt.Sprintf(`--%s=5000`, highMessageCountThresholdFlag),
	)
	if cloudProvider == cloudProviderAWS {
		args = append(
			args,
			"--iam",
			fmt.Sprintf(`--%s=pager_action`, queueAlertAlarmActionsFlag),
			fmt.Sprintf(`--%s=pager_action2`, queueAlertAlarmActionsFlag),
			fmt.Sprintf(`--%s=pager_action`, queueAlertOKActionsFlag),
			fmt.Sprintf(`--%s=pager_action2`, queueAlertOKActionsFlag),
			fmt.Sprintf(`--%s=pager_action`, dlqAlertAlarmActionsFlag),
			fmt.Sprintf(`--%s=pager_action2`, dlqAlertAlarmActionsFlag),
			fmt.Sprintf(`--%s=pager_action`, dlqAlertOKActionsFlag),
			fmt.Sprintf(`--%s=pager_action2`, dlqAlertOKActionsFlag),
		)
	} else if cloudProvider == cloudProviderGoogle {
		args = append(
			args,
			fmt.Sprintf(`--%s=alerting-project`, googleProjectAlertingFlag),
			fmt.Sprintf(`--%s=gs://myBucket/tmp`, dataflowTmpGCSLocationFlag),
			fmt.Sprintf(
				`--%s=gs://dataflow-templates/2019-04-03-00/Cloud_PubSub_to_Cloud_PubSub`,
				dataflowPubSubToPubSubTemplateGCSPathFlag,
			),
			fmt.Sprintf(
				`--%s=projects/myProject/notificationChannels/10357685029951383687`,
				queueAlertNotificationChannelsFlag,
			),
			fmt.Sprintf(
				`--%s=projects/myProject/notificationChannels/95138368710357685029`, queueAlertNotificationChannelsFlag),
			fmt.Sprintf(
				`--%s=projects/myProject/notificationChannels/10357685029951383687`, dlqAlertNotificationChannelsFlag),
			fmt.Sprintf(
				`--%s=projects/myProject/notificationChannels/95138368710357685029`, dlqAlertNotificationChannelsFlag),
			fmt.Sprintf(
				`--%s=projects/myProject/notificationChannels/10357685029951383687`, dataflowAlertNotificationChannelsFlag),
			fmt.Sprintf(
				`--%s=projects/myProject/notificationChannels/95138368710357685029`, dataflowAlertNotificationChannelsFlag),
			fmt.Sprintf(`--%s=us-west2-a`, googleDataflowZoneFlag),
			fmt.Sprintf(`--%s=us-west1`, googleDataflowRegionFlag),
			fmt.Sprintf(
				`--%s=gs://dataflow-templates/2019-04-03-00/Cloud_PubSub_to_GCS_Text`,
				dataflowPubSubToStorageGCSPathFlag,
			),
			fmt.Sprintf(`--%s=gs://myBucket/taskhawkBackup/`, googleFirehoseDataflowOutputDirectoryFlag),
			fmt.Sprintf(`--%s`, enableFirehoseAllMessages),
		)
	}
	return args
}

func TestGenerate(t *testing.T) {
	info, err := ioutil.ReadDir("test_fixtures")
	require.NoError(t, err)

	dmp := diffmatchpatch.New()

	for _, cloudDir := range info {
		if !cloudDir.IsDir() {
			continue
		}

		cloudProvider := cloudDir.Name()

		info, err := ioutil.ReadDir(filepath.Join("test_fixtures", cloudDir.Name()))
		require.NoError(t, err)

		for _, testDir := range info {
			testDirFullPath := filepath.Join("test_fixtures", cloudProvider, testDir.Name())

			infoTestDir, err := ioutil.ReadDir(testDirFullPath)
			require.NoError(t, err)

			err = os.RemoveAll("taskhawk")
			require.NoError(t, err)

			fmt.Printf("Test: %s/%s\n", cloudProvider, testDir.Name())

			configFilepath := filepath.Join(testDirFullPath, "test_config.json")

			var args []string
			if strings.Contains(testDir.Name(), "no_optional_param") {
				args = argsForTestNoOptional(cloudProvider, configFilepath)
			} else {
				args = argsForTest(cloudProvider, configFilepath)
			}

			assert.NoError(t, runApp(args))

			info, err := ioutil.ReadDir("taskhawk")
			assert.NoError(t, err)

			files := make([]string, len(info))
			for i, f := range info {
				files[i] = f.Name()
			}

			var testFiles []string
			for _, testOutputFile := range infoTestDir {
				if filepath.Ext(testOutputFile.Name()) != ".tf" {
					continue
				}
				testFiles = append(testFiles, testOutputFile.Name())
			}

			assert.Equal(t, testFiles, files)

			for _, testOutputFile := range infoTestDir {
				if filepath.Ext(testOutputFile.Name()) != ".tf" {
					continue
				}
				testOutputFileName := filepath.Join(testDirFullPath, testOutputFile.Name())
				expectedBytes, err := ioutil.ReadFile(testOutputFileName)
				require.NoError(t, err)

				expected := string(expectedBytes)

				// poor template engine
				expected = strings.Replace(
					expected, "{{tfDoNotEditStamp}}", tfDoNotEditStamp, -1)
				expected = strings.Replace(
					expected, "{{TFAWSQueueModuleVersion}}", TFAWSQueueModuleVersion, -1)
				expected = strings.Replace(
					expected, "{{TFAWSSchedulerModuleVersion}}", TFAWSSchedulerModuleVersion, -1)
				expected = strings.Replace(
					expected, "{{TFAWSLambdaModuleVersion}}", TFAWSLambdaModuleVersion, -1)
				expected = strings.Replace(
					expected, "{{TFGoogleQueueModuleVersion}}", TFGoogleQueueModuleVersion, -1)
				expected = strings.Replace(
					expected, "{{GENERATOR_VERSION}}", VERSION, -1)

				actualB, err := ioutil.ReadFile(filepath.Join("taskhawk", testOutputFile.Name()))
				require.NoError(t, err)

				assert.Equal(
					t, expected, string(actualB),
					dmp.DiffPrettyText(dmp.DiffMain(expected, string(actualB), true)),
				)
			}

			if t.Failed() {
				// so we can investigate what went wrong
				break
			}
		}

		if t.Failed() {
			// so we can investigate what went wrong
			break
		}
	}
}
