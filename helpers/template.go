package helpers

import (
	"embed"
	"strings"
	"text/template"
)

//go:embed templates/*
var templates embed.FS

func RenderTemplate[T any](tmplText string, data T) (string, error) {
	var out strings.Builder
	tmpl, _ := template.New("temp").Parse(tmplText)

	err := tmpl.Execute(&out, data)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func RenderTemplateFromFile[T any](tmplFile string, outputFile string, data T) error {
	tmplText, err := templates.ReadFile(tmplFile)
	if err != nil {
		return err
	}

	output, err := RenderTemplate(string(tmplText), data)
	if err != nil {
		return err
	}

	err = WriteTextFile(outputFile, output)
	return err
}
