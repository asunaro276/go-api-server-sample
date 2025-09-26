#!/bin/bash

# coverage.sh - Script to generate code coverage reports
# This script runs all tests with coverage analysis and generates reports

set -e

echo "Running code coverage analysis..."

# Create coverage directory if it doesn't exist
mkdir -p coverage

# Clean up previous coverage files
rm -f coverage/*.out coverage/*.html

# Run all tests with coverage
echo "Running tests with coverage..."
go test -race -coverprofile=coverage/coverage.out -covermode=atomic ./...

# Check if coverage file was generated
if [ ! -f coverage/coverage.out ]; then
    echo "Error: Coverage file not generated"
    exit 1
fi

# Generate HTML coverage report
echo "Generating HTML coverage report..."
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# Generate coverage summary
echo "Generating coverage summary..."
go tool cover -func=coverage/coverage.out > coverage/coverage.txt

# Display coverage summary
echo ""
echo "Coverage Summary:"
echo "=================="
cat coverage/coverage.txt

# Extract total coverage percentage
TOTAL_COVERAGE=$(go tool cover -func=coverage/coverage.out | grep total | awk '{print $3}')
echo ""
echo "Total Coverage: $TOTAL_COVERAGE"

# Check if coverage meets minimum threshold (80%)
COVERAGE_NUM=$(echo $TOTAL_COVERAGE | sed 's/%//')
THRESHOLD=80

if (( $(echo "$COVERAGE_NUM < $THRESHOLD" | bc -l) )); then
    echo "Warning: Coverage ($TOTAL_COVERAGE) is below threshold ($THRESHOLD%)"
    exit 1
else
    echo "Coverage ($TOTAL_COVERAGE) meets threshold ($THRESHOLD%)"
fi

# Generate detailed package coverage
echo ""
echo "Package Coverage Details:"
echo "========================="
go tool cover -func=coverage/coverage.out | grep -v total | sort -k3 -nr

echo ""
echo "Coverage analysis completed!"
echo "HTML report generated: coverage/coverage.html"
echo "Text report generated: coverage/coverage.txt"
echo "Raw coverage data: coverage/coverage.out"

# Optional: Open HTML report in browser (uncomment if desired)
# if command -v xdg-open &> /dev/null; then
#     xdg-open coverage/coverage.html
# elif command -v open &> /dev/null; then
#     open coverage/coverage.html
# fi