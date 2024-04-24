{{- define "_env_app_migrator" }}
- name: VAULT_APP_DATA_PATH
  value: {{ pluck .Values.global.env .Values.migrator.vault.data_path | first | default .Values.migrator.vault.data_path._default | join "," | quote }}

- name: VAULT_AUTH_TOKEN
  valueFrom:
    secretKeyRef:
      name: bc-wallet-tron-hdwallet
      key: vault_migrator_user_token
      optional: false
{{- end }}