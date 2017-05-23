#!/bin/sh
#
# Process the seed file.
gunzip /seed.json.gz
python /processJson.py

# Ensure mongo is up.  If this is the first time run, then the MongoDB server may take
# a few seconds to create the initial wiredtiger storage bits and the import will fail.
RETCODE=1
ATTEMPTS=0
while [ $RETCODE -ne 0 -a $ATTEMPTS -lt 20 ] ; do
    sleep 5
    nc -z -v -w5 mongodb 27017
    RETCODE=$?
    ATTEMPTS=$(($ATTEMPTS+1))
done

if [ $RETCODE -ne 0 ] ; then
    echo "Unable to connect to Mongo, exiting"
    exit $RETCODE
fi

# Drop existing data
mongo mongodb/urlinfo --eval 'db.urls.drop()'

# Import new data
mongoimport --host mongodb --db urlinfo --collection urls --type json --file /seed_output.json
