{{/*
Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
License: MIT NON-AI
*/}}

{{- define "_env_app_api" }}
- name: VAULT_APP_DATA_PATH
  value: {{ pluck .Values.global.env .Values.api.vault.data_path | first | default .Values.api.vault.data_path._default | join "," | quote }}

- name: VAULT_AUTH_TOKEN
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet
      key: vault_api_user_token
      optional: false

{{- if pluck .Values.global.env .Values.api.startupProbe.enabled | first | default .Values.api.startupProbe.enabled._default }}
- name: HEALTH_CHECK_STARTUP_ENABLED
  value: {{ pluck .Values.global.env .Values.api.startupProbe.enabled | first | default .Values.api.startupProbe.enabled._default | quote }}
- name: HEALTH_CHECK_STARTUP_HTTP_PORT
  value: {{ .Values.api.startupProbe.podSettings.httpGet.port }}
- name: HEALTH_CHECK_STARTUP_HTTP_READ_TIMEOUT
  value: {{ pluck .Values.global.env .Values.api.startupProbe.http_server.read_timeout | first | default .Values.api.startupProbe.http_server.read_timeout._default | quote }}
- name: HEALTH_CHECK_STARTUP_HTTP_WRITE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.api.startupProbe.http_server.write_timeout | first | default .Values.api.startupProbe.http_server.write_timeout._default | quote }}
- name: HEALTH_CHECK_STARTUP_HTTP_PATH
  value: {{ .Values.api.startupProbe.podSettings.path }}
{{- end }}

{{- if pluck .Values.global.env .Values.api.readinessProbe.enabled | first | default .Values.api.readinessProbe.enabled._default }}
- name: HEALTH_CHECK_READINESS_ENABLED
  value: {{ pluck .Values.global.env .Values.api.readinessProbe.enabled | first | default .Values.api.readinessProbe.enabled._default | quote }}
- name: HEALTH_CHECK_READINESS_HTTP_PORT
  value: {{ .Values.api.readinessProbe.podSettings.httpGet.port }}
- name: HEALTH_CHECK_READINESS_HTTP_READ_TIMEOUT
  value: {{ pluck .Values.global.env .Values.api.readinessProbe.http_server.read_timeout | first | default .Values.api.readinessProbe.http_server.read_timeout._default | quote }}
- name: HEALTH_CHECK_READINESS_HTTP_WRITE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.api.readinessProbe.http_server.write_timeout | first | default .Values.api.readinessProbe.http_server.write_timeout._default | quote }}
- name: HEALTH_CHECK_READINESS_HTTP_PATH
  value: {{ .Values.api.readinessProbe.podSettings.httpGet.path }}
{{- end }}

{{- if pluck .Values.global.env .Values.api.livenessProbe.enabled | first | default .Values.api.livenessProbe.enabled._default }}
- name: HEALTH_CHECK_LIVENESS_ENABLED
  value: {{ pluck .Values.global.env .Values.api.livenessProbe.enabled | first | default .Values.api.livenessProbe.enabled._default | quote }}
- name: HEALTH_CHECK_LIVENESS_HTTP_PORT
  value: {{ .Values.api.livenessProbe.podSettings.httpGet.port }}
- name: HEALTH_CHECK_LIVENESS_HTTP_READ_TIMEOUT
  value: {{ pluck .Values.global.env .Values.api.livenessProbe.http_server.read_timeout | first | default .Values.api.livenessProbe.http_server.read_timeout._default | quote }}
- name: HEALTH_CHECK_LIVENESS_HTTP_WRITE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.api.livenessProbe.http_server.write_timeout | first | default .Values.api.livenessProbe.http_server.write_timeout._default | quote }}
- name: HEALTH_CHECK_LIVENESS_HTTP_PATH
  value: {{ .Values.api.livenessProbe.podSettings.httpGet.path }}
{{- end }}

{{- if pluck .Values.global.env .Values.api.profiler.enabled | first | default .Values.api.profiler.enabled._default }}
- name: PROFILER_ENABLED
  value: {{ pluck .Values.global.env .Values.api.profiler.enabled | first | default .Values.api.profiler.enabled._default | quote }}
- name: PROFILER_HTTP_HOST
  value: {{ pluck .Values.global.env .Values.api.profiler.host | first | default .Values.api.profiler.host._default | quote }}
- name: PROFILER_HTTP_PORT
  value: {{ pluck .Values.global.env .Values.api.profiler.port | first | default .Values.api.profiler.port._default | quote }}
- name: PROFILER_HTTP_READ_TIMEOUT
  value: {{ pluck .Values.global.env .Values.api.profiler.read_timeout | first | default .Values.api.profiler.read_timeout._default | quote }}
- name: PROFILER_HTTP_WRITE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.api.profiler.write_timeout | first | default .Values.api.profiler.write_timeout._default | quote }}
- name: PROFILER_HTTP_INDEX_PATH
  value: {{ pluck .Values.global.env .Values.api.profiler.http_index_path | first | default .Values.api.profiler.http_index_path._default | quote }}
- name: PROFILER_HTTP_CMD_LINE_PATH
  value: {{ pluck .Values.global.env .Values.api.profiler.http_cmdline_path | first | default .Values.api.profiler.http_cmdline_path._default | quote }}
- name: PROFILER_HTTP_PROFILE_PATH
  value: {{ pluck .Values.global.env .Values.api.profiler.http_profile_path | first | default .Values.api.profiler.http_profile_path._default | quote }}
- name: PROFILER_HTTP_SYMBOL_PATH
  value: {{ pluck .Values.global.env .Values.api.profiler.http_symbol_path | first | default .Values.api.profiler.http_symbol_path._default | quote }}
- name: PROFILER_HTTP_TRACE_PATH
  value: {{ pluck .Values.global.env .Values.api.profiler.http_trace_path | first | default .Values.api.profiler.http_trace_path._default | quote }}
{{- end }}

- name: HDWALLET_WORDS_COUNT
  value: {{ pluck .Values.global.env .Values.api.hdwallet.words_count | first | default .Values.api.hdwallet.words_count._default | quote }}
- name: HDWALLET_CHAIN_ID
  value: {{ pluck .Values.global.env .Values.api.hdwallet.chain_id | first | default .Values.api.hdwallet.chain_id._default | quote }}
- name: HDWALLET_COIN_TYPE
  value: {{ pluck .Values.global.env .Values.api.hdwallet.coin_type | first | default .Values.api.hdwallet.coin_type._default | quote }}
- name: HDWALLET_PLUGIN_PATH
  value: {{ pluck .Values.global.env .Values.api.hdwallet.plugin_path | first | default .Values.api.hdwallet.plugin_path._default | quote }}
- name: HDWALLET_UNIX_SOCKET_DIR_PATH
  value: {{ pluck .Values.global.env .Values.common.unix_socket.dir_path | first | default .Values.common.unix_socket.dir_path._default | quote }}
- name: HDWALLET_UNIX_SOCKET_FILE_TEMPLATE
  value: {{ pluck .Values.global.env .Values.common.unix_socket.file_pattern | first | default .Values.common.unix_socket.file_pattern._default | quote }}
{{- end }}