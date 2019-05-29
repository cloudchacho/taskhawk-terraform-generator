package main

import (
	"path/filepath"
	"text/template"

	"gopkg.in/urfave/cli.v1"
)

type awsConfigWriter struct {
	c      *cli.Context
	config *AWSConfig
}

func (w *awsConfigWriter) shouldSkipFile(file string) bool {
	hasScheduleQueueApp, hasScheduleLambdaApp := false, false
	for _, app := range w.config.LambdaApps {
		if len(app.Schedule) > 0 {
			hasScheduleLambdaApp = true
			break
		}
	}
	for _, app := range w.config.QueueApps {
		if len(app.Schedule) > 0 {
			hasScheduleQueueApp = true
			break
		}
	}
	switch file {
	case schedulerFile:
		return !hasScheduleQueueApp && !hasScheduleLambdaApp
	case variablesFile:
		return !hasScheduleQueueApp
	default:
		return false
	}
}

func (w *awsConfigWriter) initTemplates() (*template.Template, error) {
	actions := map[string][]string{
		"QueueAlertAlarmActions": w.c.StringSlice(queueAlertAlarmActionsFlag),
		"QueueAlertOKActions":    w.c.StringSlice(queueAlertOKActionsFlag),
		"DLQAlertAlarmActions":   w.c.StringSlice(dlqAlertAlarmActionsFlag),
		"DLQAlertOKActions":      w.c.StringSlice(dlqAlertOKActionsFlag),
	}
	variables := map[string]string{
		"AwsAccountID": w.c.String(awsAccountIDFlag),
		"AwsRegion":    w.c.String(awsRegionFlag),
	}

	files := []string{
		appsTmplFile,
		schedulerTmplFile,
		variablesTmplFile,
	}
	templates := template.New(files[0]) // need an arbitrary name
	templates = templates.Funcs(template.FuncMap{
		"iam":                       func() bool { return w.c.Bool(iamFlag) },
		"highMessageCountThreshold": func() int { return w.c.Int(highMessageCountThresholdFlag) },
		"actions":                   func() map[string][]string { return actions },
		"variables":                 func() map[string]string { return variables },
		"hclvalue":                  hclvalue,
		"hclident":                  hclident,
		"tfDoNotEditStamp":          func() string { return tfDoNotEditStamp },
		"alerting":                  func() bool { return w.c.Bool(alertingFlag) },

		"TFAWSQueueModuleVersion":     func() string { return TFAWSQueueModuleVersion },
		"TFAWSLambdaModuleVersion":    func() string { return TFAWSLambdaModuleVersion },
		"TFAWSSchedulerModuleVersion": func() string { return TFAWSSchedulerModuleVersion },
	})

	for _, name := range files {
		_, err := templates.New(name).Parse(string(MustAsset(filepath.Join("aws", name))))
		if err != nil {
			return nil, err
		}
	}

	return templates, nil
}
