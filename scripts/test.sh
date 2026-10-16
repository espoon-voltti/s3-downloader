#!/bin/bash
# SPDX-FileCopyrightText: 2018-2020 City of Espoo
#
# SPDX-License-Identifier: MIT

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}")/.."

echo "Cleaning up previous docker containers and volumes..."
docker compose down -v

echo ""
echo "Running unit tests..."
docker build --target test -t s3-downloader:test .

echo ""
echo "Starting LocalStack..."
docker compose up -d localstack

echo ""
echo "Waiting for LocalStack to be healthy..."

ELAPSED=0
until docker compose exec localstack curl -f http://localhost:4566/_localstack/health 2>/dev/null; do
  sleep 1
  ELAPSED=$((ELAPSED + 1))
  if [ $ELAPSED -ge 60 ]; then
    echo "Timeout waiting for LocalStack to be ready"
    exit 1
  fi
done

echo ""
echo "Setting up test bucket and data..."
docker compose exec localstack sh -c '
  export AWS_ACCESS_KEY_ID=test
  export AWS_SECRET_ACCESS_KEY=test
  export AWS_DEFAULT_REGION=us-east-1
  
  echo "Document 1" > /tmp/doc1.txt &&
  echo "Document 2" > /tmp/doc2.txt &&
  echo "Nested document" > /tmp/nested.txt &&
  
  # Create bucket and upload to S3
  aws --endpoint-url=http://localhost:4566 s3 mb s3://test-bucket &&
  aws --endpoint-url=http://localhost:4566 s3 cp /tmp/doc1.txt s3://test-bucket/documents/doc1.txt &&
  aws --endpoint-url=http://localhost:4566 s3 cp /tmp/doc2.txt s3://test-bucket/documents/doc2.txt &&
  aws --endpoint-url=http://localhost:4566 s3 cp /tmp/nested.txt s3://test-bucket/documents/subfolder/nested.txt &&

  # One more than a single ListObjects page holds
  mkdir -p /tmp/many &&
  for i in $(seq 1 1001); do echo "$i" > "/tmp/many/$i.txt"; done &&
  aws --endpoint-url=http://localhost:4566 s3 cp --recursive --quiet /tmp/many s3://test-bucket/many/ &&

  echo "Test data setup complete!"
'

echo ""
echo "Running s3-downloader container..."
docker compose run --rm --build s3-downloader

echo ""
echo "Checking downloaded files in volume..."
docker run --rm -v s3-downloader_downloads:/downloads alpine sh -c '
  echo "Contents of /downloads:"
  find /downloads -type f || echo "No files found"
  ls -la /downloads/ || echo "Directory listing failed"
  echo ""
  
  # Verify expected files exist (prefix is stripped, so files are at root of /downloads)
  MISSING=0
  for file in "/downloads/doc1.txt" "/downloads/doc2.txt" "/downloads/subfolder/nested.txt"; do
    if [ -f "$file" ]; then
      echo "✓ Found: $file"
      cat "$file"
    else
      echo "✗ Missing: $file"
      MISSING=$((MISSING + 1))
    fi
  done
  
  echo ""
  if [ $MISSING -eq 0 ]; then
    echo "✓ All files downloaded successfully!"
    exit 0
  else
    echo "✗ $MISSING file(s) missing"
    exit 1
  fi
'

echo ""
echo "Downloading more objects than one listing page holds..."
docker compose run --rm --entrypoint sh s3-downloader -c '
  ./s3-downloader test-bucket many /tmp/many 2> /tmp/log || { cat /tmp/log; exit 1; }
  COUNT=$(find /tmp/many -type f | wc -l)
  if [ "$COUNT" -eq 1001 ]; then
    echo "✓ All 1001 objects downloaded"
  else
    echo "✗ Only $COUNT of 1001 objects downloaded"
    exit 1
  fi
'
