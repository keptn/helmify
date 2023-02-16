package probes

import (
	"testing"

	"github.com/arttor/helmify/pkg/helmify"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func Test_setProbeField(t *testing.T) {

	v := make(helmify.Values)
	tests := []struct {
		name     string
		deplName string
		c        *corev1.Container
		values   *helmify.Values
		wanted   string
		wantErr  bool
	}{
		{
			name:     "probe exists, values are added",
			deplName: "testdep",
			c: &corev1.Container{
				Name: "mycontainer",
				LivenessProbe: &corev1.Probe{
					Handler:        corev1.Handler{},
					TimeoutSeconds: 14,
					PeriodSeconds:  10,
				},
			},
			values:  &v,
			wanted:  "livenessProbe:\n  periodSeconds: 10\n  timeoutSeconds: 14\n",
			wantErr: false,
		},
		{
			name:     "readinessprobe exists, values are added",
			deplName: "testdep",
			c: &corev1.Container{
				Name: "mycontainer",
				ReadinessProbe: &corev1.Probe{
					Handler:        corev1.Handler{},
					TimeoutSeconds: 1,
					PeriodSeconds:  20,
				},
			},
			values:  &v,
			wanted:  "readinessProbe:\n  periodSeconds: 20\n  timeoutSeconds: 1\n",
			wantErr: false,
		},
		{
			name:     "probe is nill, no values are added",
			deplName: "testdep",
			c: &corev1.Container{
				Name: "mycontainer",
			},
			values:  &v,
			wanted:  "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v = make(helmify.Values)
			if err := setProbeField(tt.deplName, tt.c, tt.values); (err != nil) != tt.wantErr {
				t.Errorf("setProbeField() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wanted != "" {
				val := (*tt.values)["testdep"].(map[string]interface{})["mycontainer"]
				b, err := yaml.Marshal(val)
				require.Nil(t, err)
				require.Contains(t, string(b), tt.wanted)
			} else {
				require.Empty(t, (tt.values))
			}

		})
	}
}

func Test_setProbesTemplates(t *testing.T) {

	strContainers := make([]interface{}, 3)
	strContainers[1] = "previouscontent"
	strContainers[2] = "other"
	tests := []struct {
		name           string
		deploymentName string
		container      *corev1.Container
		strContainers  *interface{}
		want           string
		wantErr        bool
	}{
		{
			name:           "no probe",
			deploymentName: "test",
			container: &corev1.Container{
				Name:           "myc",
				ReadinessProbe: nil,
			},
			strContainers: &strContainers[0],
			want:          "",
			wantErr:       false,
		},
		{
			name:           "readiness probe",
			deploymentName: "test",
			container: &corev1.Container{
				Name: "myc",
				ReadinessProbe: &corev1.Probe{
					InitialDelaySeconds: 15,
				},
			},
			strContainers: &strContainers[1],

			want: "previouscontent\n" +
				"{{- if .Values.test.myc.readinessProbe }}\n" +
				"readinessProbe: {{- include \"tplvalues.render\" (dict \"value\" .Values.test.myc.readinessProbe \"context\" $) | nindent 10 }}\n" +
				" {{- else }}\n" +
				"readinessProbe:\n " +
				"initialDelaySeconds: 15\n" +
				"{{- end }}",
			wantErr: false,
		},
		{
			name:           "add liveness probe",
			deploymentName: "test",
			container: &corev1.Container{
				Name: "myc",
				LivenessProbe: &corev1.Probe{
					InitialDelaySeconds: 1,
				},
			},
			strContainers: &strContainers[2],

			want: "other\n{{- if .Values.test.myc.livenessProbe }}\n" +
				"livenessProbe: {{- include \"tplvalues.render\" (dict \"value\" .Values.test.myc.livenessProbe \"context\" $) | nindent 10 }}\n" +
				" {{- else }}\n" +
				"livenessProbe:\n initialDelaySeconds: 1\n{{- end }}",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setProbesTemplates(tt.deploymentName, tt.container, tt.strContainers, tt.container.Name); (err != nil) != tt.wantErr {
				t.Errorf("setProbesTemplates() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != "" {
				t.Log(*tt.strContainers)
				require.Contains(t, *tt.strContainers, tt.want)
			} else {
				require.Empty(t, tt.strContainers)
			}
		})
	}
}
