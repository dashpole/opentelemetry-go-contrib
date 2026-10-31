# Feature Gates

The runtime package contains feature gates used to enable the opt-in metrics of
the [OpenTelemetry Go Runtime conventions].

Note that the new runtime metrics conventions are still experimental, and may
change in backwards incompatible ways as feedback is applied.

## Features

- [Include Opt-In Metrics](#include-opt-in-metrics)

### Include Opt-In Metrics

To produce the [opt-in metrics] of the Go runtime conventions, list them
comma-separated:

```console
export OTEL_GO_X_RUNTIME_METRICS_OPTIN=go.memory.gc.pause.duration
```

The supported metrics are:

- `go.memory.gc.pause.duration`, produced by `NewProducer`.

Unknown names are ignored. The opt-in metrics produced by `Start`
(`go.memory.gc.cycles` and `go.cpu.time`) instead use the opt-in and stability
advisory parameters, and are enabled via a View.

[opt-in metrics]: https://github.com/open-telemetry/semantic-conventions/blob/main/docs/general/metric-requirement-level.md#opt-in
[OpenTelemetry Go Runtime conventions]: https://github.com/open-telemetry/semantic-conventions/blob/main/docs/runtime/go-metrics.md
