package constraints

import (
	"github.com/arttor/helmify/pkg/helmify"
	corev1 "k8s.io/api/core/v1"
)

const topologyExpression = "\n{{- if .Values.topologySpreadConstraints }}\n" +
	"      topologySpreadConstraints: {{- include \"tplvalues.render\" (dict \"value\" .Values.topologySpreadConstraints \"context\" $) | nindent 8 }}\n" +
	"{{- end }}\n"

const nodeSelectorExpression = "{{- if .Values.nodeSelector }}\n" +
	"      nodeSelector: {{- include \"tplvalues.render\" ( dict \"value\" .Values.nodeSelector \"context\" $) | nindent 8 }}\n" +
	"{{- end }}\n"

const tolerationsExpression = "{{- if .Values.tolerations }}\n" +
	"      tolerations: {{- include \"tplvalues.render\" (dict \"value\" .Values.tolerations \"context\" .) | nindent 8 }}\n" +
	"{{- end }}\n"

// ProcessSpecMap adds 'topologyConstraints' to the podSpec in specMap, if it doesn't
// already has one defined.
func ProcessSpecMap(specMap string, values *helmify.Values, podspec corev1.PodSpec) string {

	(*values)["topologySpreadConstraints"] = podspec.TopologySpreadConstraints
	(*values)["nodeSelector"] = podspec.NodeSelector
	(*values)["tolerations"] = podspec.Tolerations

	topology := (*values)["topologySpreadConstraints"].([]corev1.TopologySpreadConstraint)
	if len(topology) == 0 {
		(*values)["topologySpreadConstraints"] = []interface{}{}
	}
	nodeSelector := (*values)["nodeSelector"].(map[string]string)
	if len(nodeSelector) == 0 {
		(*values)["nodeSelector"] = map[string]string{}
	}

	tolerations := (*values)["tolerations"].([]corev1.Toleration)

	if len(tolerations) == 0 {
		(*values)["tolerations"] = []interface{}{}
	}

	return specMap + topologyExpression + nodeSelectorExpression + tolerationsExpression
}
