#!/usr/bin/env bash

# Import Graph Analysis Tool for SDK
# Detects circular dependencies and verifies dependency direction
# Usage: ./scripts/check_imports.sh [--json]

set -euo pipefail

SDK_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUTPUT_JSON=false

# Parse arguments
if [[ "${1:-}" == "--json" ]]; then
    OUTPUT_JSON=true
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check for Go
if ! command -v go &> /dev/null; then
    echo "Error: go command not found"
    exit 1
fi

cd "$SDK_ROOT/go"

# Generate import graph
echo "Analyzing import graph..." >&2
GRAPH=$(go list -f '{{.ImportPath}}: {{join .Imports ","}}' ./... 2>&1)

# Check for circular dependencies
echo "Checking for circular dependencies..." >&2
CIRCULAR_DEPS=$(go mod graph | grep "reglet-sdk.*reglet-sdk" | grep -v "^github.com/reglet-dev/reglet-sdk" || true)

if [ -n "$CIRCULAR_DEPS" ]; then
    if $OUTPUT_JSON; then
        echo '{"status":"fail","error":"circular_dependencies","details":"'"${CIRCULAR_DEPS}"'"}'
    else
        echo -e "${RED}❌ CIRCULAR DEPENDENCIES DETECTED:${NC}"
        echo "$CIRCULAR_DEPS"
    fi
    exit 1
fi

# Verify hexagonal architecture dependency direction
# Rule: infrastructure -> application -> domain
# Domain should never import from application or infrastructure
echo "Checking hexagonal architecture compliance..." >&2

DOMAIN_IMPORTS=$(echo "$GRAPH" | grep "domain/" | grep -E "(infrastructure|application)" || true)

if [ -n "$DOMAIN_IMPORTS" ]; then
    if $OUTPUT_JSON; then
        echo '{"status":"fail","error":"architecture_violation","details":"Domain layer imports from infrastructure/application","violations":"'"${DOMAIN_IMPORTS}"'"}'
    else
        echo -e "${RED}❌ ARCHITECTURE VIOLATION:${NC}"
        echo "Domain layer should not import from application or infrastructure layers"
        echo "$DOMAIN_IMPORTS"
    fi
    exit 1
fi

# Application should not import from infrastructure
APPLICATION_IMPORTS=$(echo "$GRAPH" | grep "application/" | grep "infrastructure" || true)

if [ -n "$APPLICATION_IMPORTS" ]; then
    if $OUTPUT_JSON; then
        echo '{"status":"fail","error":"architecture_violation","details":"Application layer imports from infrastructure","violations":"'"${APPLICATION_IMPORTS}"'"}'
    else
        echo -e "${YELLOW}⚠️  WARNING:${NC}"
        echo "Application layer should avoid importing from infrastructure layer (prefer ports)"
        echo "$APPLICATION_IMPORTS"
    fi
fi

# Check that SDK doesn't import Reglet
SDK_REGLET_IMPORTS=$(echo "$GRAPH" | grep -v "^reglet/" | grep "reglet/" || true)

if [ -n "$SDK_REGLET_IMPORTS" ]; then
    if $OUTPUT_JSON; then
        echo '{"status":"fail","error":"forbidden_dependency","details":"SDK imports from Reglet (should be reversed)","violations":"'"${SDK_REGLET_IMPORTS}"'"}'
    else
        echo -e "${RED}❌ FORBIDDEN DEPENDENCY:${NC}"
        echo "SDK must not import Reglet (dependency should be Reglet -> SDK)"
        echo "$SDK_REGLET_IMPORTS"
    fi
    exit 1
fi

# All checks passed
if $OUTPUT_JSON; then
    echo '{"status":"pass","message":"All import graph checks passed"}'
else
    echo -e "${GREEN}✅ All import graph checks passed${NC}"
    echo "  - No circular dependencies"
    echo "  - Hexagonal architecture compliant"
    echo "  - SDK does not import Reglet"
fi

exit 0
