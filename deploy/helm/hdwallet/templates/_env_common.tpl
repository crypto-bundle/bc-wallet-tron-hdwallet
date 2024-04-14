{{- define "_env_common" }}
- name: APP_ENV
  value: {{ pluck .Values.global.env .Values.common.environment | first | default .Values.common.environment._default | quote }}
- name: APP_DEBUG
  value: {{ pluck .Values.global.env .Values.common.debug | first | default .Values.common.debug._default | quote }}
- name: APP_STAGE
  value: {{ pluck .Values.global.env .Values.common.stage | first | default .Values.common.stage._default | quote }}

- name: LOGGER_LEVEL
  value: {{ pluck .Values.global.env .Values.common.logger.minimal_level | first | default .Values.common.logger.minimal_level._default | quote }}
- name: LOGGER_STACKTRACE_ENABLE
  value: {{ pluck .Values.global.env .Values.common.logger.enabled_stack_trace | first | default .Values.common.logger.enabled_stack_trace._default | quote }}

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
{{- end }}