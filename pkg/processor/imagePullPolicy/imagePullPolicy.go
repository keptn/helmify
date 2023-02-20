package imagePullPolicy

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/keptn/helmify/pkg/helmify"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const helmTemplate = "{{ .Values.%[1]s.%[2]s.imagePullPolicy }}"

// ProcessSpecMap adds 'imagePullSecrets' to the podSpec in specMap, if it doesn't
// already has one defined.
func ProcessSpecMap(name string, specMap map[string]interface{}, values *helmify.Values) error {

	cs, _, err := unstructured.NestedSlice(specMap, "containers")

	if err != nil {
		return err
	}

	newContainers := make([]interface{}, len(cs))
	for i, c := range cs {
		castedContainer := c.(map[string]interface{})
		containerName := strcase.ToLowerCamel(castedContainer["name"].(string))
		if castedContainer["imagePullPolicy"] != nil {
			err = unstructured.SetNestedField(*values, castedContainer["imagePullPolicy"], name, containerName, "imagePullPolicy")
			if err != nil {
				return err
			}
			castedContainer["imagePullPolicy"] = fmt.Sprintf(helmTemplate, name, containerName)
		}
		newContainers[i] = castedContainer
	}
	return unstructured.SetNestedSlice(specMap, newContainers, "containers")
}
