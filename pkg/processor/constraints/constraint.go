package constraints

import (
	"github.com/keptn/helmify/pkg/helmify"
	corev1 "k8s.io/api/core/v1"
)

const tolerations = "tolerations"
const topology = "topologySpreadConstraints"
const nodeSelector = "nodeSelector"

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

	(*values)[topology] = podspec.TopologySpreadConstraints
	(*values)[nodeSelector] = podspec.NodeSelector
	(*values)[tolerations] = podspec.Tolerations

	tp := (*values)[topology].([]corev1.TopologySpreadConstraint)
	if len(tp) == 0 {
		(*values)[topology] = []interface{}{}
	}
	ns := (*values)[nodeSelector].(map[string]string)
	if len(ns) == 0 {
		(*values)[nodeSelector] = map[string]string{}
	}

	tl := (*values)[tolerations].([]corev1.Toleration)

	if len(tl) == 0 {
		(*values)[tolerations] = []interface{}{}
	}

	return specMap + topologyExpression + nodeSelectorExpression + tolerationsExpression
}
