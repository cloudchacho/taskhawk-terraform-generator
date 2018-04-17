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
	alertsFile    = "alerts.tf"
	appsFile      = "apps.tf"
	schedulerFile = "scheduler.tf"
)

func shouldSkipFile(config *Config, needAlerts bool, file string) bool {
	switch file {
	case alertsFile:
		return !needAlerts
	case schedulerFile:
		hasSchedule := false
		for _, app := range config.LambdaApps {
			if len(app.Schedule) > 0 {
				hasSchedule = true
				break
			}
		}
		if hasSchedule {
			return false
		}
		for _, app := range config.QueueApps {
			if len(app.Schedule) > 0 {
				hasSchedule = true
				break
			}
		}
		return !hasSchedule
	default:
		return false
	}
}

func writeFiles(config *Config, needAlerts bool, module string, templates *template.Template) error {
	files := []string{alertsFile, appsFile, schedulerFile}

	for _, file := range files {
		if shouldSkipFile(config, needAlerts, file) {
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

	templates := template.New("taskhawk-templates")
	templates = templates.Funcs(template.FuncMap{
		"generator_version": func() string { return VERSION },
		"version":           func() string { return TFModulesVersion },
		"iam":               func() bool { return c.Bool(iamFlag) },
		"actions":           func() map[string][]string { return actions },
		"hclvalue":          hclvalue,
		"hclident":          hclident,
		"tfDoNotEditStamp":  func() string { return tfDoNotEditStamp },
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

	needAlerts := c.Bool(alertingFlag) && len(config.QueueApps) > 0

	if err := writeFiles(config, needAlerts, module, templates); err != nil {
		return err
	}

	if err := hclFmtDir(module); err != nil {
		return err
	}

	return nil
}
