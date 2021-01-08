package naming

import (
	"bytes"
	"html/template"
	"strings"
)

var tplFuncs = map[string]interface{}{
	"upper": strings.ToUpper,
	"title": strings.Title,
	"lower": strings.ToLower,
}

type NamingTemplate struct {
	template       *template.Template
	data           map[string]interface{}
	namingTemplate string
}

func NewNamingTemplate(namingTemplate string, data map[string]interface{}) (*NamingTemplate, error) {
	t, err := template.New("template").Funcs(tplFuncs).Parse(namingTemplate)
	if err != nil {
		return nil, err
	}
	return &NamingTemplate{
		template:       t,
		namingTemplate: namingTemplate,
		data:           data,
	}, nil
}

func (n *NamingTemplate) GetBindingName(bindingName string) (string, error) {
	d := map[string]interface{}{
		"service": n.data,
		"name":    bindingName,
	}

	var tpl bytes.Buffer
	err := n.template.Execute(&tpl, d)
	if err != nil {
		return bindingName, err
	}
	return tpl.String(), nil
}
