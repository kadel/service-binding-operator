package naming

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildNamingStrategy(t *testing.T) {
	data := map[string]interface{}{
		"name": "database",
		"kind": "Service",
	}

	dataProvider := []struct {
		namingTemplate      string
		bindingName         string
		expectedBindingName string
		expectedError       error
	}{
		{
			namingTemplate:      "{{ .service.kind | upper }}",
			bindingName:         "db",
			expectedBindingName: "SERVICE",
		},
		{
			namingTemplate:      "{{ .service.kind | upper }}_{{ .name }}",
			bindingName:         "db",
			expectedBindingName: "SERVICE_db",
		},
		{
			namingTemplate:      "{{ .service.kind | upper }}_{{ .name | upper }}",
			bindingName:         "db",
			expectedBindingName: "SERVICE_DB",
		},
		{
			namingTemplate:      "{{ .wrongfield | upper }}",
			bindingName:         "db",
			expectedBindingName: "db",
			expectedError:       errors.New("template: template:1:17: executing \"template\" at <upper>: invalid value; expected string"),
		},
	}

	for _, tt := range dataProvider {
		t.Run(fmt.Sprintf("Test with template %s", tt.namingTemplate), func(t *testing.T) {
			template, err := NewNamingTemplate(tt.namingTemplate, data)
			assert.NoError(t, err)
			bindingName, err := template.GetBindingName(tt.bindingName)
			if err != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			}
			assert.EqualValues(t, tt.expectedBindingName, bindingName)
		})
	}

}
