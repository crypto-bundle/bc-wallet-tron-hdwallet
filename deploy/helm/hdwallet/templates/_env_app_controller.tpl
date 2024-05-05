{{- define "_env_app_controller" }}
- name: VAULT_APP_DATA_PATH
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

- name: HDWALLET_UNIX_SOCKET_PATH
  value: {{ pluck .Values.global.env .Values.common.unix_socket_path | first | default .Values.common.unix_socket_path._default | quote }}

{{- end }}