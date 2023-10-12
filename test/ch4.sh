echo -n "this object will have only 1 instance" | openssl dgst -sha256 -binary | base64

curl -v 10.29.2.1:12345/objects/test4_1 -XPUT -d "this object will have only 1 instance" -H "Digest: SHA-256=aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY="
curl 10.29.2.1:12345/locate/aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY=
ls /tmp/?/objects/aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY\=