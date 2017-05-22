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

    if port is None:
        if parsedUrl.scheme == "https":
            port = 443
        else:
            port = 80

    pathquery = parsedUrl.path
    if parsedUrl.query is not None and parsedUrl.query != "":
        pathquery += "?%s" % parsedUrl.query

    return {
        'hostport' : "%s:%d" % (parsedUrl.netloc, port),
        'pathquery': pathquery
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
