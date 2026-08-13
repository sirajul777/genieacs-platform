# Extension SDK

The Agent exposes a small stable Go contract for GenieACS integrations.

## Implement an extension

Implement `extension.Extension`:

- `Name()` returns the unique extension name.
- `Execute(ctx, request)` performs the operation.
- `Request.Payload` is opaque JSON/bytes owned by the extension.
- `Response.Payload` contains the extension result.

Example:

```go
type MyExtension struct{}

func (MyExtension) Name() string { return "my-extension" }

func (MyExtension) Execute(ctx context.Context, req extension.Request) (extension.Response, error) {
    return extension.Response{Payload: req.Payload}, nil
}
```

Register implementations with `extension.NewRegistry(...)` and execute them through `extension.NewRunner(...)`.

The runner propagates cancellation and rejects unknown extension names. The core Agent does not interpret extension payloads.
