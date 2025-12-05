{{- define "myapp.name" -}}
{{ .Chart.Name }}
{{- end -}}

{{- define "myapp.fullname" -}}
{{ .Release.Name }}-{{ .Chart.Name }}
{{- end -}}

