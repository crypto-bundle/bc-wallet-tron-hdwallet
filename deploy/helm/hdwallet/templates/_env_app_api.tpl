{{- define "_env_app_api" }}
- name: VAULT_APP_DATA_PATH
  value: {{ pluck .Values.global.env .Values.app.vault.data_path | first | default .Values.app.vault.data_path._default | join "," | quote }}

- name: VAULT_AUTH_TOKEN
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet-api
      key: vault_api_user_token
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


- name: HDWALLET_WORDS_COUNT
  value: {{ pluck .Values.global.env .Values.app.mnemonic.words_count | first | default .Values.app.mnemonic.words_count._default | quote }}
- name: HDWALLET_PLUGIN_PATH
  value: {{ pluck .Values.global.env .Values.app.mnemonic.plugin_path | first | default .Values.app.mnemonic.plugin_path._default | quote }}
- name: HDWALLET_UNIX_SOCKET_PATH
  value: {{ pluck .Values.global.env .Values.common.unix_socket_path | first | default .Values.common.unix_socket_path._default | quote }}
{{- end }}