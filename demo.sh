#!/bin/bash

# Demo script for Binks CLI Stage 1
echo "=== Binks CLI Stage 1 Demo ==="
echo

echo "1. Testing basic echo command:"
./binks echo "Hello, DevShell!"
echo

echo "2. Testing file listing:"
./binks ls -la
echo

echo "3. Testing current directory:"
./binks pwd
echo

echo "4. Testing shell features (environment variables):"
./binks echo "User: \$USER, Home: \$HOME"
echo

echo "5. Testing argument handling:"
./binks echo hello world
echo
./binks echo "hello, world"
echo
./binks echo 'Hello, DevShell!'
echo

echo "6. Testing multiline output:"
./binks find . -name "*.go" | head -5
echo

echo "7. Testing error handling with invalid command:"
./binks invalidcommand123 2>&1 || echo "Error handling works!"
echo

echo "8. Testing usage message:"
./binks 2>&1 || echo "Usage message displayed"
echo

echo "=== Demo Complete ==="
