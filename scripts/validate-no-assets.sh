#!/bin/bash
# validate-no-assets.sh
# Validates that no pre-bundled assets exist in the repository.
# This enforces the "100% procedural generation" constraint.

set -e

PROHIBITED_EXTENSIONS="png jpg jpeg gif svg bmp ico webp mp3 wav ogg flac aac m4a wma"

echo "Checking for prohibited asset files..."

found_any=false
for ext in $PROHIBITED_EXTENSIONS; do
    # Find files with this extension, excluding vendor and .git directories
    files=$(find . -type f -name "*.$ext" \
        -not -path "./.git/*" \
        -not -path "./vendor/*" \
        -not -path "./.github/*" \
        2>/dev/null || true)
    
    if [ -n "$files" ]; then
        echo "FAIL: Found prohibited .$ext files:"
        echo "$files"
        found_any=true
    fi
done

if [ "$found_any" = true ]; then
    echo ""
    echo "ERROR: Pre-bundled assets are not allowed in this project."
    echo "All visual and audio assets must be procedurally generated at runtime."
    exit 1
fi

echo "PASS: No bundled assets found"
exit 0
