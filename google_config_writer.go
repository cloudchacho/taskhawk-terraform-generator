package main

import (
	"gopkg.in/urfave/cli.v1"
	"path/filepath"
	"text/template"
)

type googleConfigWriter struct {
	c      *cli.Context
	config *GoogleConfig
}

func (w *googleConfigWriter) shouldSkipFile(file string) bool {
	return file == schedulerFile
}

func (w *googleConfigWriter) initTemplates() (*template.Template, error) {
	channels := map[string][]string{
		"QueueAlertNotificationChannels": w.c.StringSlice(queueAlertNotificationChannelsFlag),
		"DLQAlertNotificationChannels":   w.c.StringSlice(dlqAlertNotificationChannelsFlag),
	}
	variables := map[string]string{
		"DataflowTmpGCSLocation":                 w.c.String(dataflowTmpGCSLocationFlag),
		"DataflowPubSubToPubSubTemplateGCSPath":  w.c.String(dataflowPubSubToPubSubTemplateGCSPathFlag),
		"DataflowPubSubToStorageTemplateGCSPath": w.c.String(dataflowPubSubToStorageGCSPathFlag),
		"DataflowZone":                           w.c.String(googleDataflowZoneFlag),
		"DataflowRegion":                           w.c.String(googleDataflowRegionFlag),
		"DataflowOutputDirectory":                w.c.String(googleFirehoseDataflowOutputDirectoryFlag),
		"GoogleProjectAlerting":                  w.c.String(googleProjectAlertingFlag),
	}
	flags := map[string]bool{
		"EnableFirehoseAllMessages": w.c.Bool(enableFirehoseAllMessages),
	}
	files := []string{
		appsTmplFile,
		variablesTmplFile,
	}
	templates := template.New(files[0]) // need an arbitrary name
	templates = templates.Funcs(template.FuncMap{
		"highMessageCountThreshold": func() int { return w.c.Int(highMessageCountThresholdFlag) },
		"channels":                  func() map[string][]string { return channels },
		"variables":                 func() map[string]string { return variables },
		"flags":                     func() map[string]bool { return flags },
		"hclvalue":                  hclvalue,
		"hclident":                  hclident,
		"tfDoNotEditStamp":          func() string { return tfDoNotEditStamp },
		"alerting":                  func() bool { return w.c.Bool(alertingFlag) },

		"TFGoogleQueueModuleVersion": func() string { return TFGoogleQueueModuleVersion },
	})

	for _, name := range files {
		_, err := templates.New(name).Parse(string(MustAsset(filepath.Join("google", name))))
		if err != nil {
			return nil, err
		}
	}

	return templates, nil
}
