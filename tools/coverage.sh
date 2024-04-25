#!/bin/sh -ue

VIEWER=${VIEWER:-xdg-open}

COVER_PROFILE="artifacts/coverage.out"

REPORT_TEXT="artifacts/coverage_report_last.txt"

PKG="./..."

COVER_MODE=${COVER_MODE:-"set"}

REPORT_HTML="artifacts/coverage.html"

mkdir -p "./artifacts"

go test -timeout=10s -coverpkg="$PKG" -coverprofile="$COVER_PROFILE" -covermode="$COVER_MODE" ./...

if [ "${1:-""}" = html ]; then
    go tool cover -html="$COVER_PROFILE" -o "$REPORT_HTML"
    exec $VIEWER "$REPORT_HTML"
fi

go tool cover -func="$COVER_PROFILE" | tee "$REPORT_TEXT"
