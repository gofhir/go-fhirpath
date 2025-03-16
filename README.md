# fhirpath-evaluator

A lightweight Go module to evaluate **FHIRPath expressions** against FHIR resource JSON.

## Features
- 🩺 Full FHIRPath evaluation using ANTLR parser.
- 🔌 Pluggable handlers (e.g., `where()`, `exists()`, `empty()`)
- 🔄 Hot-reloadable configuration (enable/disable handlers at runtime)
- ✅ Modular and production-ready

## Install

```bash
go get github.com/yourname/fhirpath-evaluator