{{- define "_env_store" }}
- name: POSTGRESQL_SERVICE_HOST
  value: {{ pluck .Values.global.env .Values.common.db.host | first | default .Values.common.db.host._default | quote }}
- name: POSTGRESQL_SERVICE_PORT
  value: {{ pluck .Values.global.env .Values.common.db.port | first | default .Values.common.db.port._default | quote }}
- name: POSTGRESQL_SSL_MODE
  value: {{ pluck .Values.global.env .Values.common.db.ssl_mode | first | default .Values.common.db.ssl_mode._default | quote }}
- name: POSTGRESQL_MAX_OPEN_CONNECTIONS
  value: {{ pluck .Values.global.env .Values.common.db.open_connections | first | default .Values.common.db.open_connections._default | quote }}
- name: POSTGRESQL_MAX_IDLE_CONNECTIONS
  value: {{ pluck .Values.global.env .Values.common.db.idle_connections | first | default .Values.common.db.idle_connections._default | quote }}
- name: POSTGRESQL_CONNECTION_RETRY_COUNT
  value: {{ pluck .Values.global.env .Values.common.db.connection_retry_count | first | default .Values.common.db.connection_retry_count._default | quote }}
- name: POSTGRESQL_CONNECTION_RETRY_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.db.connection_retry_timeout | first | default .Values.common.db.connection_retry_timeout._default | quote }}

- name: REDIS_HOST
  value: {{ pluck .Values.global.env .Values.common.redis.host | first | default .Values.common.redis.host._default | quote }}
- name: REDIS_PORT
  value: {{ pluck .Values.global.env .Values.common.redis.port | first | default .Values.common.redis.port._default | quote }}
- name: REDIS_DB
  value: {{ pluck .Values.global.env .Values.common.redis.db | first | default .Values.common.redis.db._default | quote }}
- name: REDIS_CONNECTION_RETRY_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.redis.connection_retry_timeout | first | default .Values.common.redis.connection_retry_timeout._default | quote }}
- name: REDIS_CONNECTION_RETRY_COUNT
  value: {{ pluck .Values.global.env .Values.common.redis.connection_retry_count | first | default .Values.common.redis.connection_retry_count._default | quote }}
- name: REDIS_MAX_RETRY_COUNT
  value: {{ pluck .Values.global.env .Values.common.redis.max_retry_count | first | default .Values.common.redis.max_retry_count._default | quote }}
- name: REDIS_READ_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.redis.read_timeout | first | default .Values.common.redis.read_timeout._default | quote }}
- name: REDIS_WRITE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.redis.write_timeout | first | default .Values.common.redis.write_timeout._default | quote }}
- name: REDIS_MIN_IDLE_CONNECTIONS
  value: {{ pluck .Values.global.env .Values.common.redis.min_idle_connections | first | default .Values.common.redis.min_idle_connections._default | quote }}
- name: REDIS_IDLE_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.redis.idle_timeout | first | default .Values.common.redis.idle_timeout._default | quote }}
- name: REDIS_MAX_CONNECTION_AGE
  value: {{ pluck .Values.global.env .Values.common.redis.connection_age | first | default .Values.common.redis.connection_age._default | quote }}
- name: REDIS_POOL_SIZE
  value: {{ pluck .Values.global.env .Values.common.redis.pool_size | first | default .Values.common.redis.pool_size._default | quote }}
- name: REDIS_POOL_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.redis.pool_timeout | first | default .Values.common.redis.pool_timeout._default | quote }}
- name: REDIS_DIAL_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.redis.dial_timeout | first | default .Values.common.redis.dial_timeout._default | quote }}

- name: NATS_ADDRESSES
  value: {{ pluck .Values.global.env .Values.common.nats.hosts | first | default .Values.common.nats.hosts._default | join "," | quote }}
- name: NATS_CONNECTION_RETRY
  value: {{ pluck .Values.global.env .Values.common.nats.connection_retry | first | default .Values.common.nats.connection_retry._default | quote }}
- name: NATS_CONNECTION_RETRY_COUNT
  value: {{ pluck .Values.global.env .Values.common.nats.connection_retry_count | first | default .Values.common.nats.connection_retry_count._default | quote }}
- name: NATS_CONNECTION_RETRY_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.nats.connection_retry_timeout | first | default .Values.common.nats.connection_retry_timeout._default | quote }}
- name: NATS_FLUSH_TIMEOUT
  value: {{ pluck .Values.global.env .Values.common.nats.flush_timeout | first | default .Values.common.nats.flush_timeout._default | quote }}
- name: NATS_WORKER_PER_CONSUMER
  value: {{ pluck .Values.global.env .Values.common.nats.workers | first | default .Values.common.nats.workers._default | quote }}
- name: NATS_KV_BUCKET_REPLICAS
  value: {{ pluck .Values.global.env .Values.common.nats.kv_bucket_replicas | first | default .Values.common.nats.kv_bucket_replicas._default | quote }}
{{- end }}
