#!/usr/bin/env bash

cd presentation && go test -bench=BenchmarkGenerateViews && goimports -w presentation/web_view.GEN.go
