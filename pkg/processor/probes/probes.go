package probes

import (
	"fmt"

	"github.com/arttor/helmify/pkg/helmify"
	yamlformat "github.com/arttor/helmify/pkg/yaml"
	"github.com/iancoleman/strcase"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

const livenessProbe = "\n{{- if .Values.%[1]s.%[2]s.livenessProbe }}\n" +
	"livenessProbe: {{- include \"tplvalues.render\" (dict \"value\" .Values.%[1]s.%[2]s.livenessProbe \"context\" $) | nindent 10 }}\n" +
	" {{- else }}\n" +
	"livenessProbe:\n%[3]s" +
	"\n{{- end }}"

const readinessProbe = "\n{{- if .Values.%[1]s.%[2]s.readinessProbe }}\n" +
	"readinessProbe: {{- include \"tplvalues.render\" (dict \"value\" .Values.%[1]s.%[2]s.readinessProbe \"context\" $) | nindent 10 }}\n" +
	" {{- else }}\n" +
	"readinessProbe:\n%[3]s" +
	"\n{{- end }}"

// ProcessSpecMap adds 'probes' to the Containers in specMap, if they are defined
func ProcessSpecMap(name string, specMap map[string]interface{}, values *helmify.Values, pspec corev1.PodSpec) error {

	strContainers := make([]interface{}, len(pspec.Containers))
	cs, _, err := unstructured.NestedSlice(specMap, "containers")
	if err != nil {
		return err
	}

	for i := range cs {
		containerName := strcase.ToLowerCamel(pspec.Containers[i].Name)
		var ready, live string
		content, err := yamlformat.Marshal(cs[i], 0)
		if err != nil {
			return err
		}
		strContainers[i] = content
		if pspec.Containers[i].LivenessProbe != nil {

			live, err = yamlformat.Marshal(pspec.Containers[i].LivenessProbe, 1)
			if err != nil {
				return err
			}
			strContainers[i] = strContainers[i].(string) +
				fmt.Sprintf(livenessProbe, name, containerName, live)
		}
		if pspec.Containers[i].ReadinessProbe != nil {

			ready, err = yamlformat.Marshal(pspec.Containers[i].ReadinessProbe, 1)
			if err != nil {
				return err
			}
			strContainers[i] = strContainers[i].(string) +
				fmt.Sprintf(readinessProbe, name, containerName, ready)
		}
		setProbeField(name, pspec.Containers[i], values)
	}
	unstructured.SetNestedSlice(specMap, strContainers, "containers")

	return err
}

func setProbeField(name string, c corev1.Container, values *helmify.Values) corev1.Container {

	containerName := strcase.ToLowerCamel(c.Name)
	if c.LivenessProbe != nil {
		ready, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(c.LivenessProbe)
		unstructured.SetNestedField(*values, ready, name, containerName, "livenessProbe")
	}
	if c.ReadinessProbe != nil {
		ready, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(c.ReadinessProbe)
		unstructured.SetNestedField(*values, ready, name, containerName, "readinessProbe")
	}
	return c
}
