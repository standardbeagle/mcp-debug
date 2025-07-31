#!/bin/bash
# Master test runner for MCP Debug

set -e

echo "=== MCP Debug Test Suite ==="
echo

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Track test results
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run a test
run_test() {
    local test_name=$1
    local test_command=$2
    
    echo -n "Running $test_name... "
    
    if eval "$test_command" > /tmp/test_output.log 2>&1; then
        echo -e "${GREEN}PASSED${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}FAILED${NC}"
        echo "  Error output:"
        tail -n 20 /tmp/test_output.log | sed 's/^/    /'
        ((TESTS_FAILED++))
    fi
}

# Build main binary
echo "Building mcp-debug binary..."
cd ..
go build -o mcp-debug .
cd tests

# Build test servers
echo "Building test servers..."
cd ../test-servers
for server in math-server file-server lifecycle-server-v1 lifecycle-server-v2; do
    echo "  Building $server..."
    go build -o $server ${server}.go
done
cd ../tests

echo
echo "Running Integration Tests:"
echo "=========================="

# Run Python integration tests
cd integration
run_test "Proxy Calls Test" "python3 test-proxy-calls.py"
run_test "Dynamic Registration Test" "python3 test-dynamic-registration.py"
run_test "Lifecycle Test" "python3 test-lifecycle.py"
run_test "Simple Dynamic Test" "python3 test-simple-dynamic.py"
run_test "Updated Tools Test" "python3 test-updated-tools.py"
cd ..

echo
echo "Running Script Tests:"
echo "===================="

# Run shell script tests
cd scripts
run_test "Playback Test" "bash test-playback.sh"
cd ..

echo
echo "Test Summary:"
echo "============="
echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "\n${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "\n${RED}Some tests failed!${NC}"
    exit 1
fi