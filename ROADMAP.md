# TODO
## Websockets

- Add field that would represent shard identifier
  - Make it possible to filter by shard
  - Do not pass data to log processor if another shard is selected or excluded
## Storage
- Add Grafana Loki support
  - Implement Query Builder

## Security

- Develop authentication and authorization for Opensearch Dashboards requests
- Design shard access control
  - Make possible to handle several indexes by one shard
  - Create authorization rules for shards (JWT/NoAuth)
  