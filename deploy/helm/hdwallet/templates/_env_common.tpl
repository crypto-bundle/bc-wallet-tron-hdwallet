{{ define "_env_common" }}
- name: APP_ENV
  value: {{ pluck .Values.global.env .Values.common.environment | first | default .Values.common.environment._default | quote }}
- name: APP_DEBUG
  value: {{ pluck .Values.global.env .Values.common.debug_mode | first | default .Values.common.debug_mode._default | quote }}
- name: APP_STAGE
  value: {{ pluck .Values.global.env .Values.common.stage.name | first | default .Values.common.stage.name._default | quote }}

- name: LOGGER_LEVEL
  value: {{ pluck .Values.global.env .Values.common.logger.minimal_level | first | default .Values.common.logger.minimal_level._default | quote }}
- name: LOGGER_STACKTRACE_ENABLE
  value: {{ pluck .Values.global.env .Values.common.logger.enabled_stack_trace | first | default .Values.common.logger.enabled_stack_trace._default | quote }}

- name: GOMEMLIMIT
  valueFrom:
    resourceFieldRef:
      resource: limits.memory

- name: GOMAXPROCS
  valueFrom:
    resourceFieldRef:
      resource: limits.cpu
{{ end }}