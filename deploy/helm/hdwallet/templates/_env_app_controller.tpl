{{- define "_env_app_controller" }}
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
{{- end }}