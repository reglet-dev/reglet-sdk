#!/usr/bin/env bash

# Benchmark Comparison Tool for SDK
# Compares current benchmarks against baseline to detect ≤5% regression
# Usage: ./scripts/benchmark_compare.sh [baseline_file]

set -euo pipefail

SDK_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BASELINE_FILE="${1:-benchmarks/baseline.txt}"
CURRENT_FILE="/tmp/current_benchmarks.txt"
REGRESSION_THRESHOLD=5.0  # 5% regression threshold

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check for Go
if ! command -v go &> /dev/null; then
    echo "Error: go command not found"
    exit 1
fi

# Check for benchstat (optional but recommended)
HAS_BENCHSTAT=false
if command -v benchstat &> /dev/null; then
    HAS_BENCHSTAT=true
fi

cd "$SDK_ROOT/go"

# Run current benchmarks
echo "Running current benchmarks..." >&2
go test ./... -bench=. -benchmem -run=^$ > "$CURRENT_FILE" 2>&1 || true

if [ ! -f "$BASELINE_FILE" ]; then
    echo -e "${YELLOW}⚠️  No baseline file found at $BASELINE_FILE${NC}"
    echo "Creating baseline from current benchmarks..."
    mkdir -p "$(dirname "$BASELINE_FILE")"
    cp "$CURRENT_FILE" "$BASELINE_FILE"
    echo -e "${GREEN}✅ Baseline created${NC}"
    exit 0
fi

# Compare with baseline
echo "Comparing with baseline..." >&2

if $HAS_BENCHSTAT; then
    # Use benchstat for detailed comparison
    COMPARISON=$(benchstat "$BASELINE_FILE" "$CURRENT_FILE" 2>&1)
    echo "$COMPARISON"
    
    # Check for significant regressions
    # benchstat shows ~X% changes; extract percentages and check threshold
    REGRESSIONS=$(echo "$COMPARISON" | grep -E '\+[0-9]+\.[0-9]+%' | awk '{
        # Extract percentage from columns like "+10.5%"
        match($0, /\+([0-9]+\.[0-9]+)%/, arr);
        if (arr[1] > '"$REGRESSION_THRESHOLD"') {
            print $0
        }
    }')
    
    if [ -n "$REGRESSIONS" ]; then
        echo -e "${RED}❌ PERFORMANCE REGRESSION DETECTED (>${REGRESSION_THRESHOLD}%):${NC}"
        echo "$REGRESSIONS"
        exit 1
    fi
    
    echo -e "${GREEN}✅ No significant performance regressions detected${NC}"
else
    # Simple comparison without benchstat
    echo -e "${YELLOW}⚠️  benchstat not found. Install with: go install golang.org/x/perf/cmd/benchstat@latest${NC}"
    echo "Performing basic comparison..."
    
    # Extract ns/op values and compare
    while IFS= read -r line; do
        if [[ "$line" =~ ^Benchmark ]]; then
            BENCH_NAME=$(echo "$line" | awk '{print $1}')
            CURRENT_NSOP=$(echo "$line" | awk '{for(i=1;i<=NF;i++) if($i ~ /ns\/op/) print $(i-1)}')
            
            # Find baseline for same benchmark
            BASELINE_LINE=$(grep "^$BENCH_NAME" "$BASELINE_FILE" || true)
            if [ -n "$BASELINE_LINE" ]; then
                BASELINE_NSOP=$(echo "$BASELINE_LINE" | awk '{for(i=1;i<=NF;i++) if($i ~ /ns\/op/) print $(i-1)}')
                
                if [ -n "$CURRENT_NSOP" ] && [ -n "$BASELINE_NSOP" ]; then
                    # Calculate percentage change
                    PERCENT_CHANGE=$(awk -v curr="$CURRENT_NSOP" -v base="$BASELINE_NSOP" 'BEGIN {
                        if (base > 0) {
                            change = ((curr - base) / base) * 100;
                            printf "%.2f", change
                        } else {
                            print "0"
                        }
                    }')
                    
                    # Check threshold
                    IS_REGRESSION=$(awk -v change="$PERCENT_CHANGE" -v threshold="$REGRESSION_THRESHOLD" 'BEGIN {
                        if (change > threshold) print "1"; else print "0"
                    }')
                    
                    if [ "$IS_REGRESSION" == "1" ]; then
                        echo -e "${RED}❌ REGRESSION:${NC} $BENCH_NAME: $BASELINE_NSOP -> $CURRENT_NSOP ns/op (+${PERCENT_CHANGE}%)"
                        exit 1
                    fi
                fi
            fi
        fi
    done < "$CURRENT_FILE"
    
    echo -e "${GREEN}✅ No significant performance regressions detected${NC}"
fi

exit 0
