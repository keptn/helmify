package security_context

import (
	"testing"

	"github.com/arttor/helmify/pkg/helmify"
	"github.com/stretchr/testify/assert"
)

//func TestProcessContainerSecurityContext(t *testing.T) {
//	type args struct {
//		nameCamel string
//		specMap   map[string]interface{}
//		values    *helmify.Values
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ProcessContainerSecurityContext(tt.args.nameCamel, tt.args.specMap, tt.args.values)
//		})
//	}
//}

func Test_setSecContextValue(t *testing.T) {
	type args struct {
		resourceName            string
		containerName           string
		castedContainer         map[string]interface{}
		values                  *helmify.Values
		fieldName               string
		useRenderedHelmTemplate bool
	}
	tests := []struct {
		name string
		args args
		want *helmify.Values
	}{
		{
			name: "test if value is generated correctly",
			args: args{
				resourceName:  "someResource",
				containerName: "someContainer",
				castedContainer: map[string]interface{}{
					"securityContext": map[string]interface{}{
						"someField": "someValue",
					},
				},
				values:                  &helmify.Values{},
				fieldName:               "someField",
				useRenderedHelmTemplate: false,
			},
			want: &helmify.Values{
				"someResource": map[string]interface{}{
					"someContainer": map[string]interface{}{
						"containerSecurityContext": map[string]interface{}{
							"someField": "someValue",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setSecContextValue(tt.args.resourceName, tt.args.containerName, tt.args.castedContainer, tt.args.values, tt.args.fieldName, tt.args.useRenderedHelmTemplate)
			assert.Equal(t, tt.want, tt.args.values)
		})
	}
}
