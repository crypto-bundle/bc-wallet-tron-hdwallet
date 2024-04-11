{{- define "_env_app_api" }}
- name: VAULT_DATA_PATH
  value: {{ pluck .Values.global.env .Values.app.vault.data_path | first | default .Values.app.vault.data_path._default | join "," | quote }}

- name: VAULT_AUTH_TOKEN
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet-api
      key: vault_api_user_token
      optional: false

- name: VAULT_TRANSIT_KEY
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet
      key: vault_transit_secret_key_api
      optional: false

- name: HDWALLET_WORDS_COUNT
  value: {{ pluck .Values.global.env .Values.app.mnemonic.words_count | first | default .Values.app.mnemonic.words_count._default | quote }}
{{- end }}