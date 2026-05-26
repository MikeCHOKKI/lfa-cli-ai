#!/usr/bin/env bash
set -euo pipefail

BINARY="${1:-./bin/lfa}"
PASS=0
FAIL=0

green() { printf "\033[32m%s\033[0m\n" "$*"; }
red()   { printf "\033[31m%s\033[0m\n" "$*"; }

test_name()   { printf "\n=== %s ===\n" "$1"; }
pass()        { green "  PASS"; PASS=$((PASS + 1)); }
fail()        { red "  FAIL: $1"; FAIL=$((FAIL + 1)); }

test_name "doctor (should run without errors)"
if output=$("$BINARY" doctor 2>&1); then
  pass
else
  fail "$output"
fi

test_name "version (should print version)"
output=$("$BINARY" version 2>&1)
if [[ "$output" == *"v0.1.0"* ]]; then
  pass
else
  fail "expected v0.1.0, got: $output"
fi

test_name "setup --dry-run (should simulate)"
output=$("$BINARY" setup --dry-run 2>&1)
if [[ "$output" == *"Dry-run"* ]]; then
  pass
else
  fail "expected dry-run message, got: $output"
fi

test_name "dashboard --yes (non-interactive)"
output=$("$BINARY" dashboard --yes 2>&1)
if echo "$output" | grep -qE "(Setup complete|installed|error)"; then
  pass
else
  fail "unexpected output: $output"
fi

printf "\n=========================="
printf "\n  %d passed, %d failed\n" "$PASS" "$FAIL"
printf "==========================\n"

[[ "$FAIL" -eq 0 ]]
