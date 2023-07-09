# Logger

[![GoDoc](https://godoc.org/github.com/GuanceCloud/logger-go?status.svg)](https://godoc.org/github.com/GuanceCloud/logger-go)
[![MIT License](https://img.shields.io/badge/license-MIT-green?style=plastic)](LICENSE)

This is a logger wrap on [zap](https://github.com/uber-go/zap) for simple daily usage, and aimed to replace [exist logger](https://github.com/GuanceCloud/cliutils/tree/main/logger).

This repo still on development and do not promise API compatible.

## Done

See *zap_test.go* use cases.

## TODO

- [ ] Add `guance-syncer` to push logging data to Guance Cloud
- [ ] Add trace-id inject(use [otelzap](https://github.com/uptrace/opentelemetry-go-extra/tree/main/otelzap))
