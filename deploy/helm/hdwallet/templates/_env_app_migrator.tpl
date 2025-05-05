{{/*
Copyright (c) 2022-2025 Aleksei Kotelnikov(gudron2s@gmail.com)
License: MIT NON-AI
*/}}

{{- define "_env_app_migrator" }}
- name: VAULT_APP_DATA_PATH
  value: {{ pluck .Values.global.env .Values.migrator.vault.data_path | first | default .Values.migrator.vault.data_path._default | join "," | quote }}

- name: VAULT_AUTH_TOKEN_RENEW_TTL
  value: {{ pluck .Values.global.env .Values.migrator.vault.renew_ttl | first | default .Values.migrator.vault.renew_ttl._default | quote }}

- name: VAULT_AUTH_TOKEN_FILE_PATH
  value: {{ pluck .Values.global.env .Values.migrator.vault.token_path | first | default .Values.migrator.vault.token_path._default | quote }}

{{- end }}