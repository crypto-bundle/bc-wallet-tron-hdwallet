{{- define "_env_app_controller" }}
{{- if .Values.global.env .Values.common.healthcheck.startup.enabled | first | default .Values.common.healthcheck.startup.enabled._default }}
- name: HEALTH_CHECK_STARTUP_ENABLED
  value: {{ pluck .Values.global.env .Values.common.healthcheck.startup.enabled | first | default .Values.common.healthcheck.startup.enabled._default | quote }}
- name: HEALTH_CHECK_STARTUP_HTTP_PORT
  value: {{ .Values.common.startupProbe.httpGet.port }}
- name: HEALTH_CHECK_STARTUP_HTTP_READ_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.healthcheck.startup.http_server.read_timeout | first | default .Values.common.healthcheck.startup.http_server.read_timeout._default | quote }}
- name: HEALTH_CHECK_STARTUP_HTTP_WRITE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.healthcheck.startup.http_server.write_timeout | first | default .Values.common.healthcheck.startup.http_server.write_timeout._default | quote }}
- name: HEALTH_CHECK_STARTUP_HTTP_PATH
  value: {{ .Values.common.startupProbe.httpGet.path }}
{{- end }}

{{- if .Values.global.env .Values.common.healthcheck.readiness.enabled | first | default .Values.common.healthcheck.readiness.enabled._default }}
- name: HEALTH_CHECK_READINESS_ENABLED
  value: {{ pluck .Values.global.env .Values.common.healthcheck.readiness.enabled | first | default .Values.common.healthcheck.readiness.enabled._default | quote }}
- name: HEALTH_CHECK_READINESS_HTTP_PORT
  value: {{ .Values.common.readinessProbe.httpGet.port }}
- name: HEALTH_CHECK_READINESS_HTTP_READ_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.healthcheck.readiness.http_server.read_timeout | first | default .Values.common.healthcheck.readiness.http_server.read_timeout._default | quote }}
- name: HEALTH_CHECK_READINESS_HTTP_WRITE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.healthcheck.readiness.http_server.write_timeout | first | default .Values.common.healthcheck.readiness.http_server.write_timeout._default | quote }}
- name: HEALTH_CHECK_READINESS_HTTP_PATH
  value: {{ .Values.common.readinessProbe.httpGet.path }}
{{- end }}

{{- if .Values.global.env .Values.common.healthcheck.liveness.enabled | first | default .Values.common.healthcheck.liveness.enabled._default }}
- name: HEALTH_CHECK_LIVENESS_ENABLED
  value: {{ pluck .Values.global.env .Values.common.healthcheck.liveness.enabled | first | default .Values.common.healthcheck.liveness.enabled._default | quote }}
- name: HEALTH_CHECK_LIVENESS_HTTP_PORT
  value: {{ .Values.common.livenessProbe.httpGet.port }}
- name: HEALTH_CHECK_LIVENESS_HTTP_READ_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.healthcheck.liveness.http_server.read_timeout | first | default .Values.common.healthcheck.liveness.http_server.read_timeout._default | quote }}
- name: HEALTH_CHECK_LIVENESS_HTTP_WRITE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.healthcheck.liveness.http_server.write_timeout | first | default .Values.common.healthcheck.liveness.http_server.write_timeout._default | quote }}
- name: HEALTH_CHECK_LIVENESS_HTTP_PATH
  value: {{ .Values.common.livenessProbe.httpGet.path }}
{{- end }}

- name: VAULT_DATA_PATH
  value: {{ pluck .Values.global.env .Values.app.vault.data_path | first | default .Values.app.vault.data_path._default | join "," | quote }}

- name: VAULT_AUTH_TOKEN
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet-controller
      key: vault_controller_user_token
      optional: false

- name: VAULT_COMMON_TRANSIT_KEY
  valueFrom:
    secretKeyRef:
      name: bc-wallet-common
      key: vault_transit_secret_key
      optional: false

- name: VAULT_APP_ENCRYPTION_KEY
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet
      key: vault_transit_secret_key_controller
      optional: false

- name: API_GRPC_PORT
  value: {{ pluck .Values.global.env .Values.app.controller.grpc_port | first | default .Values.app.controller.grpc_port._default | quote }}

- name: EVENT_CHANNEL_WORKERS_COUNT
  value: {{ pluck .Values.global.env .Values.app.controller.events.workers_count | first | default .Values.app.controller.events.workers_count._default | quote }}
- name: EVENT_CHANNEL_BUFFER_SIZE
  value: {{ pluck .Values.global.env .Values.app.controller.events.buffer_size | first | default .Values.app.controller.events.buffer_size._default | quote }}

{{- end }}