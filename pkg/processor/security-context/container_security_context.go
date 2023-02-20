package security_context

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/keptn/helmify/pkg/helmify"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	sc                   = "securityContext"
	cscValueName         = "containerSecurityContext"
	allowPrivEsc         = "allowPrivilegeEscalation"
	caps                 = "capabilities"
	privileged           = "privileged"
	roFS                 = "readOnlyRootFilesystem"
	gid                  = "runAsGroup"
	nonRoot              = "runAsNonRoot"
	uid                  = "runAsUser"
	seccompProfile       = "seccompProfile"
	helmTemplate         = "{{ .Values.%[1]s.%[2]s.%[3]s.%[4]s }}"
	helmTemplateRendered = "{{- include \"tplvalues.render\" (dict \"value\" .Values.%[1]s.%[2]s.%[3]s.%[4]s \"context\" $) | nindent 12 }}"
)

// ProcessContainerSecurityContext adds 'securityContext' to the podSpec in specMap, if it doesn't have one already defined.
func ProcessContainerSecurityContext(nameCamel string, specMap map[string]interface{}, values *helmify.Values) {
	if _, defined := specMap["containers"]; defined {
		containers, _, _ := unstructured.NestedSlice(specMap, "containers")
		for _, container := range containers {
			castedContainer := container.(map[string]interface{})
			containerName := strcase.ToLowerCamel(castedContainer["name"].(string))
			if _, defined2 := castedContainer["securityContext"]; defined2 {
				setSecContextValue(nameCamel, containerName, castedContainer, values, allowPrivEsc, false)
				setSecContextValue(nameCamel, containerName, castedContainer, values, privileged, false)
				setSecContextValue(nameCamel, containerName, castedContainer, values, roFS, false)
				setSecContextValue(nameCamel, containerName, castedContainer, values, gid, false)
				setSecContextValue(nameCamel, containerName, castedContainer, values, nonRoot, false)
				setSecContextValue(nameCamel, containerName, castedContainer, values, uid, false)
				setSecContextValue(nameCamel, containerName, castedContainer, values, caps, true)
				setSecContextValue(nameCamel, containerName, castedContainer, values, seccompProfile, true)
			}
		}
		unstructured.SetNestedSlice(specMap, containers, "containers")
	}
}

func setSecContextValue(resourceName string, containerName string, castedContainer map[string]interface{}, values *helmify.Values, fieldName string, useRenderedHelmTemplate bool) {
	if castedContainer["securityContext"].(map[string]interface{})[fieldName] != nil {
		unstructured.SetNestedField(*values, castedContainer["securityContext"].(map[string]interface{})[fieldName], resourceName, containerName, cscValueName, fieldName)

		valueString := fmt.Sprintf(helmTemplate, resourceName, containerName, cscValueName, fieldName)

		if useRenderedHelmTemplate {
			valueString = fmt.Sprintf(helmTemplateRendered, resourceName, containerName, cscValueName, fieldName)
		}

		unstructured.SetNestedField(castedContainer, valueString, sc, fieldName)
	}
}
