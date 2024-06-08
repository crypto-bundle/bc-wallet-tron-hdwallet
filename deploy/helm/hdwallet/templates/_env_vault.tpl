{{/*
Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
License: MIT NON-AI
*/}}

{{- define "_env_vault" }}
- name: VAULT_SERVICE_HOST
  value: {{ pluck .Values.global.env .Values.common.vault.host | first | default .Values.common.vault.host._default | quote }}
- name: VAULT_SERVICE_PORT
  value: {{ pluck .Values.global.env .Values.common.vault.port | first | default .Values.common.vault.port._default | quote }}
- name: VAULT_USE_HTTPS
  value: {{ pluck .Values.global.env .Values.common.vault.use_https | first | default .Values.common.vault.use_https._default | quote }}
- name: VAULT_AUTH_METHOD
  value: {{ pluck .Values.global.env .Values.common.vault.auth_method | first | default .Values.common.vault.auth_method._default }}
{{- end }}