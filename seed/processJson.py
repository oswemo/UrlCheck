#!/usr/bin/env python
#
import os
import json
from urlparse import urlparse

def readSeedFile(filename):
    rawData = None
    with open(filename, "r") as fp:
        rawData = fp.read()

    return json.loads(rawData)

def processEntry(entry):
    url = entry.get('url', '')
    parsedUrl = urlparse(url)

    port = parsedUrl.port

    if port == None:
        if parsedUrl.scheme == "http":
            port = 80
        if parsedUrl.scheme == "https":
            port = 443

    return {
        'hostport' : "%s:%d" % (parsedUrl.netloc, port),
        'pathquery': "%s?%s" % (parsedUrl.path, parsedUrl.query)
    }

#############################################
filename = "/seed.json"
output = "/seed_output.json"

if not os.path.exists(filename):
    raise Exception("Missing file: %s", filename)

fp = open(output, "w")
for entry in readSeedFile(filename):
    data = processEntry(entry)
    fp.write(json.dumps(data))

fp.close()
