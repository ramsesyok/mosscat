#!/usr/bin/env bash
set -euo pipefail

npx --yes openapi-typescript-codegen \
  --input ../internal/api/openapi.yaml \
  --output src/api \
  --client axios \
  --useOptions \
  --useUnionTypes
