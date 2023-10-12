#!/bin/bash
# init
dd if=/dev/urandom of=/tmp/file bs=1000 count=100
hash=$(openssl dgst -sha256 -binary /tmp/file | base64)
dd if=/tmp/file of=/tmp/first bs=1000 count=50
dd if=/tmp/file of=/tmp/second bs=1000 skip=32 count=68
token=$(curl -v 10.29.2.1:12345/objects/test6 -XPOST -H "Digest: SHA-256=$hash" -H "Size: 100000")
echo $token
