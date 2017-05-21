FROM mongo:3.4

COPY processJson.py /
COPY import.sh /
COPY seed.json.gz /

RUN chmod +x /import.sh
RUN apt-get update && apt-get install -y python netcat

CMD /import.sh
