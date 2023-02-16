package topologyConstraint

import (
	"fmt"

	"github.com/arttor/helmify/pkg/helmify"
)

const helmExpression = "\n{{- if .Values.topologySpreadConstraints }}\n" +
	"      topologySpreadConstraints: {{- include \"tplvalues.render\" (dict \"value\" .Values.%s.topologySpreadConstraints \"context\" $) | nindent 8 }}\n" +
	"{{- end }}"

// ProcessSpecMap adds 'topologyConstraints' to the podSpec in specMap, if it doesn't
// already has one defined.
func ProcessSpecMap(name string, specMap string, values *helmify.Values) string {

	(*values)["topologyConstraints"] = []string{}
	return specMap + fmt.Sprintf(helmExpression, name)
}
