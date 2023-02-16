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

	cs, _, err := unstructured.NestedSlice(specMap, "containers")
	if err != nil {
		return err
	}

	strContainers, err := templateContainers(name, cs, pspec, values)
	if err != nil {
		return err
	}
	return unstructured.SetNestedSlice(specMap, strContainers, "containers")

}

func templateContainers(name string, cs []interface{}, pspec corev1.PodSpec, values *helmify.Values) ([]interface{}, error) {
	strContainers := make([]interface{}, len(pspec.Containers))
	for i := range cs {
		containerName := strcase.ToLowerCamel(pspec.Containers[i].Name)

		content, err := yamlformat.Marshal(cs[i], 0)
		if err != nil {
			return nil, err
		}
		strContainers[i] = content
		err = setProbesTemplates(name, &(pspec.Containers[i]), &strContainers[i], containerName)
		if err != nil {
			return nil, err
		}
		err = setProbeField(name, &(pspec.Containers[i]), values)
		if err != nil {
			return nil, err
		}
	}
	return strContainers, nil
}

func setProbesTemplates(name string, container *corev1.Container, strContainers *interface{}, containerName string) error {

	if container.LivenessProbe != nil {
		live, err := yamlformat.Marshal(container.LivenessProbe, 1)
		if err != nil {
			return err
		}
		*strContainers = (*strContainers).(string) + fmt.Sprintf(livenessProbe, name, containerName, live)
	}
	if container.ReadinessProbe != nil {
		ready, err := yamlformat.Marshal(container.ReadinessProbe, 1)
		if err != nil {
			return err
		}
		*strContainers = (*strContainers).(string) + fmt.Sprintf(readinessProbe, name, containerName, ready)
	}
	return nil

}

func setProbeField(name string, c *corev1.Container, values *helmify.Values) error {

	containerName := strcase.ToLowerCamel(c.Name)
	if c.LivenessProbe != nil {
		ready, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(c.LivenessProbe)
		err := unstructured.SetNestedField(*values, ready, name, containerName, "livenessProbe")
		if err != nil {
			return err
		}
	}
	if c.ReadinessProbe != nil {
		ready, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(c.ReadinessProbe)
		return unstructured.SetNestedField(*values, ready, name, containerName, "readinessProbe")
	}
	return nil
}