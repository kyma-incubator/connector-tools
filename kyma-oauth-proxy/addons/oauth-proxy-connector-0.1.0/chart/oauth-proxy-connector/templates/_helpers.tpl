{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "bundle.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" $name .Values.application_name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "oauth_proxy_service" -}}
{{- printf "http://%s-%s.%s:8080" .Values.oAuthProxyName .Values.application_name .Release.Namespace | trimAll " " | quote -}}
{{- end -}}

{{- define "api_spec_url_value" -}}
{{- $specUrl := default .Values.default_api_spec_url .Values.api_spec_url -}}
{{- printf "%s" $specUrl -}}
{{- end -}}


    

