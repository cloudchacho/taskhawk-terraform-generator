package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

const (
	appsFile      = "apps.tf"
	schedulerFile = "scheduler.tf"
	variablesFile = "variables.tf"
)

func shouldSkipFile(config *Config, file string) bool {
	hasScheduleQueueApp, hasScheduleLambdaApp := false, false
	for _, app := range config.LambdaApps {
		if len(app.Schedule) > 0 {
			hasScheduleLambdaApp = true
			break
		}
	}
	for _, app := range config.QueueApps {
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

func writeFiles(config *Config, module string, templates *template.Template) error {
	files := []string{appsFile, schedulerFile, variablesFile}

	for _, file := range files {
		if shouldSkipFile(config, file) {
			continue
		}
		path := filepath.Join(module, file)
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		templateName := fmt.Sprintf("%s.tmpl", file)
		err = templates.ExecuteTemplate(f, templateName, config)
		errClose := f.Close()
		if err != nil {
			return err
		}
		if errClose != nil {
			return errClose
		}
	}
	return nil
}

func initTemplates(c *cli.Context) (*template.Template, error) {
	actions := map[string][]string{
		"QueueAlertAlarmActions": c.StringSlice(queueAlertAlarmActionsFlag),
		"QueueAlertOKActions":    c.StringSlice(queueAlertOKActionsFlag),
		"DLQAlertAlarmActions":   c.StringSlice(dlqAlertAlarmActionsFlag),
		"DLQAlertOKActions":      c.StringSlice(dlqAlertOKActionsFlag),
	}
	variables := map[string]string{
		"AwsAccountID": c.String(awsAccountIDFlag),
		"AwsRegion":    c.String(awsRegionFlag),
	}

	templates := template.New("taskhawk-templates")
	templates = templates.Funcs(template.FuncMap{
		"generator_version":         func() string { return VERSION },
		"iam":                       func() bool { return c.Bool(iamFlag) },
		"actions":                   func() map[string][]string { return actions },
		"variables":                 func() map[string]string { return variables },
		"hclvalue":                  hclvalue,
		"hclident":                  hclident,
		"tfDoNotEditStamp":          func() string { return tfDoNotEditStamp },
		"alerting":                  func() bool { return c.Bool(alertingFlag) },
		"highMessageCountThreshold": func() int { return c.Int(highMessageCountThresholdFlag) },

		"TFQueueModuleVersion":     func() string { return TFQueueModuleVersion },
		"TFLambdaModuleVersion":    func() string { return TFLambdaModuleVersion },
		"TFSchedulerModuleVersion": func() string { return TFSchedulerModuleVersion },
	})

	for _, name := range AssetNames() {
		_, err := templates.New(name).Parse(string(MustAsset(name)))
		if err != nil {
			return nil, err
		}
	}

	return templates, nil
}

func writeTerraform(config *Config, c *cli.Context) error {
	module := c.String(moduleFlag)

	if len(module) == 0 {
		return fmt.Errorf("invalid module")
	}

	if err := os.Mkdir(module, os.ModePerm); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("module directory already exists: %q", module)
		}
		return fmt.Errorf("unable to create module dir: %q, error: %q", module, err)
	}

	templates, err := initTemplates(c)
	if err != nil {
		return errors.Wrap(err, "unable to initialize templates")
	}

	if err := writeFiles(config, module, templates); err != nil {
		return err
	}

	if err := hclFmtDir(module); err != nil {
		return err
	}

	return nil
}
