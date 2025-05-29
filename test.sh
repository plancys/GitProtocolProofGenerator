#!/bin/bash

# Test script for Git Report Generator

echo "🚀 Git Report Generator Test Script"
echo "=================================="

# Check if git-report-generator binary exists
if [ ! -f "./git-report-generator" ]; then
    echo "❌ Binary not found. Building..."
    make build
    if [ $? -ne 0 ]; then
        echo "❌ Build failed!"
        exit 1
    fi
fi

echo "✅ Binary found"

# Test help command
echo ""
echo "📖 Testing help command:"
echo "------------------------"
./git-report-generator --help

# Test with invalid date format
echo ""
echo "🧪 Testing invalid date format:"
echo "-------------------------------"
./git-report-generator --from invalid-date --to 2024-12-31 --author test@example.com 2>&1 | head -1

# Test with missing required flags
echo ""
echo "🧪 Testing missing required flags:"
echo "----------------------------------"
./git-report-generator 2>&1 | head -1

# Test with non-existent repository
echo ""
echo "🧪 Testing non-existent repository:"
echo "-----------------------------------"
./git-report-generator --repo /non/existent/path --from 2024-01-01 --to 2024-12-31 --author test@example.com 2>&1 | head -1

echo ""
echo "✅ All tests completed!"
echo ""
echo "📝 To test with a real repository:"
echo "  1. Navigate to a Git repository with commits"
echo "  2. Run: ./git-report-generator --from YYYY-MM-DD --to YYYY-MM-DD"
echo "  3. Or specify author: ./git-report-generator --from YYYY-MM-DD --to YYYY-MM-DD --author your@email.com" 