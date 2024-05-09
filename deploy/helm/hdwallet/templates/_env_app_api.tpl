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

- name: HDWALLET_WORDS_COUNT
  value: {{ pluck .Values.global.env .Values.api.mnemonic.words_count | first | default .Values.api.mnemonic.words_count._default | quote }}
- name: HDWALLET_PLUGIN_PATH
  value: {{ pluck .Values.global.env .Values.api.mnemonic.plugin_path | first | default .Values.api.mnemonic.plugin_path._default | quote }}
- name: HDWALLET_UNIX_SOCKET_DIR_PATH
  value: {{ pluck .Values.global.env .Values.common.unix_socket.dir_path | first | default .Values.common.unix_socket.dir_path._default | quote }}
- name: HDWALLET_UNIX_SOCKET_FILE_TEMPLATE
  value: {{ pluck .Values.global.env .Values.common.unix_socket.file_pattern | first | default .Values.common.unix_socket.file_pattern._default | quote }}
{{- end }}