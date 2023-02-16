package topologyConstraint

import (
	"github.com/arttor/helmify/pkg/helmify"
	corev1 "k8s.io/api/core/v1"
)

const helmExpression = "\n{{- if .Values.topologySpreadConstraints }}\n" +
	"      topologySpreadConstraints: {{- include \"tplvalues.render\" (dict \"value\" .Values.topologySpreadConstraints \"context\" $) | nindent 8 }}\n" +
	"{{- end }}\n"

// ProcessSpecMap adds 'topologyConstraints' to the podSpec in specMap, if it doesn't
// already has one defined.
func ProcessSpecMap(specMap string, values *helmify.Values, constraints []corev1.TopologySpreadConstraint) string {

	(*values)["topologySpreadConstraints"] = constraints
	return specMap + helmExpression
}
