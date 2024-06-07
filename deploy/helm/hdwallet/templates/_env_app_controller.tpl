{{- define "_env_app_controller" }}
- name: VAULT_APP_DATA_PATH
  value: {{ pluck .Values.global.env .Values.controller.vault.data_path | first | default .Values.controller.vault.data_path._default | join "," | quote }}

- name: VAULT_AUTH_TOKEN
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet
      key: vault_controller_user_token
      optional: false

{{- if pluck .Values.global.env .Values.controller.startupProbe.enabled | first | default .Values.controller.startupProbe.enabled._default }}
- name: HEALTH_CHECK_STARTUP_ENABLED
  value: {{ pluck .Values.global.env .Values.controller.startupProbe.enabled | first | default .Values.controller.startupProbe.enabled._default | quote }}
- name: HEALTH_CHECK_STARTUP_HTTP_PORT
  value: {{ .Values.controller.startupProbe.podSettings.httpGet.port }}
- name: HEALTH_CHECK_STARTUP_HTTP_READ_TIMEOUT
  value: {{ pluck .Values.global.env .Values.controller.startupProbe.http_server.read_timeout | first | default .Values.controller.startupProbe.http_server.read_timeout._default | quote }}
- name: HEALTH_CHECK_STARTUP_HTTP_WRITE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.controller.startupProbe.http_server.write_timeout | first | default .Values.controller.startupProbe.http_server.write_timeout._default | quote }}
- name: HEALTH_CHECK_STARTUP_HTTP_PATH
  value: {{ .Values.controller.startupProbe.podSettings.path }}
{{- end }}

{{- if pluck .Values.global.env .Values.controller.readinessProbe.enabled | first | default .Values.controller.readinessProbe.enabled._default }}
- name: HEALTH_CHECK_READINESS_ENABLED
  value: {{ pluck .Values.global.env .Values.controller.readinessProbe.enabled | first | default .Values.controller.readinessProbe.enabled._default | quote }}
- name: HEALTH_CHECK_READINESS_HTTP_PORT
  value: {{ .Values.controller.readinessProbe.podSettings.httpGet.port }}
- name: HEALTH_CHECK_READINESS_HTTP_READ_TIMEOUT
  value: {{ pluck .Values.global.env .Values.controller.readinessProbe.http_server.read_timeout | first | default .Values.controller.readinessProbe.http_server.read_timeout._default | quote }}
- name: HEALTH_CHECK_READINESS_HTTP_WRITE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.controller.readinessProbe.http_server.write_timeout | first | default .Values.controller.readinessProbe.http_server.write_timeout._default | quote }}
- name: HEALTH_CHECK_READINESS_HTTP_PATH
  value: {{ .Values.controller.readinessProbe.podSettings.httpGet.path }}
{{- end }}

{{- if pluck .Values.global.env .Values.controller.livenessProbe.enabled | first | default .Values.controller.livenessProbe.enabled._default }}
- name: HEALTH_CHECK_LIVENESS_ENABLED
  value: {{ pluck .Values.global.env .Values.controller.livenessProbe.enabled | first | default .Values.controller.livenessProbe.enabled._default | quote }}
- name: HEALTH_CHECK_LIVENESS_HTTP_PORT
  value: {{ .Values.controller.livenessProbe.podSettings.httpGet.port }}
- name: HEALTH_CHECK_LIVENESS_HTTP_READ_TIMEOUT
  value: {{ pluck .Values.global.env .Values.controller.livenessProbe.http_server.read_timeout | first | default .Values.controller.livenessProbe.http_server.read_timeout._default | quote }}
- name: HEALTH_CHECK_LIVENESS_HTTP_WRITE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.controller.livenessProbe.http_server.write_timeout | first | default .Values.controller.livenessProbe.http_server.write_timeout._default | quote }}
- name: HEALTH_CHECK_LIVENESS_HTTP_PATH
  value: {{ .Values.controller.livenessProbe.podSettings.httpGet.path }}
{{- end }}

{{- if pluck .Values.global.env .Values.controller.profiler.enabled | first | default .Values.controller.profiler.enabled._default }}
- name: PROFILER_ENABLED
  value: {{ pluck .Values.global.env .Values.controller.profiler.enabled | first | default .Values.controller.profiler.enabled._default | quote }}
- name: PROFILER_HTTP_HOST
  value: {{ pluck .Values.global.env .Values.controller.profiler.host | first | default .Values.controller.profiler.host._default | quote }}
- name: PROFILER_HTTP_PORT
  value: {{ pluck .Values.global.env .Values.controller.profiler.port | first | default .Values.controller.profiler.port._default | quote }}
- name: PROFILER_HTTP_READ_TIMEOUT
  value: {{ pluck .Values.global.env .Values.controller.profiler.read_timeout | first | default .Values.controller.profiler.read_timeout._default | quote }}
- name: PROFILER_HTTP_WRITE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.controller.profiler.write_timeout | first | default .Values.controller.profiler.write_timeout._default | quote }}
- name: PROFILER_HTTP_INDEX_PATH
  value: {{ pluck .Values.global.env .Values.controller.profiler.http_index_path | first | default .Values.controller.profiler.http_index_path._default | quote }}
- name: PROFILER_HTTP_CMD_LINE_PATH
  value: {{ pluck .Values.global.env .Values.controller.profiler.http_cmdline_path | first | default .Values.controller.profiler.http_cmdline_path._default | quote }}
- name: PROFILER_HTTP_PROFILE_PATH
  value: {{ pluck .Values.global.env .Values.controller.profiler.http_profile_path | first | default .Values.controller.profiler.http_profile_path._default | quote }}
- name: PROFILER_HTTP_SYMBOL_PATH
  value: {{ pluck .Values.global.env .Values.controller.profiler.http_symbol_path | first | default .Values.controller.profiler.http_symbol_path._default | quote }}
- name: PROFILER_HTTP_TRACE_PATH
  value: {{ pluck .Values.global.env .Values.controller.profiler.http_trace_path | first | default .Values.controller.profiler.http_trace_path._default | quote }}
{{- end }}

- name: MANAGER_API_GRPC_PORT
  value: {{ pluck .Values.global.env .Values.controller.grpc_port.manager_api | first | default .Values.controller.grpc_port.manager_api._default | quote }}

- name: WALLET_API_GRPC_PORT
  value: {{ pluck .Values.global.env .Values.controller.grpc_port.wallet_api | first | default .Values.controller.grpc_port.wallet_api._default | quote }}

- name: JWT_DEFAULT_TTL
  value: {{ pluck .Values.global.env .Values.controller.jwt.ttl | first | default .Values.controller.jwt.ttl._default | quote }}

- name: EVENT_CHANNEL_WORKERS_COUNT
  value: {{ pluck .Values.global.env .Values.controller.events.workers_count | first | default .Values.controller.events.workers_count._default | quote }}
- name: EVENT_CHANNEL_BUFFER_SIZE
  value: {{ pluck .Values.global.env .Values.controller.events.buffer_size | first | default .Values.controller.events.buffer_size._default | quote }}

- name: HDWALLET_UNIX_SOCKET_DIR_PATH
  value: {{ pluck .Values.global.env .Values.common.unix_socket.dir_path | first | default .Values.common.unix_socket.dir_path._default | quote }}
- name: HDWALLET_UNIX_SOCKET_FILE_TEMPLATE
  value: {{ pluck .Values.global.env .Values.common.unix_socket.file_pattern | first | default .Values.common.unix_socket.file_pattern._default | quote }}

{{- end }}