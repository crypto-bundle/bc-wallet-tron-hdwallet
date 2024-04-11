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
{{- end }}