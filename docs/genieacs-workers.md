# GenieACS workers

The Agent provides concrete worker types backed by the GenieACS extension:

- `genieacs.install`
- `genieacs.verify`

Job payload:

```json
{"args":{"version":"latest"}}
```

The worker reports progress through the Agent `ProgressReporter` and delegates the actual host operation to the `genieacs` extension. The extension command is configured by the Agent composition layer.
