FROM ubuntu
RUN apt-get update && apt-get install -y wget
COPY bin/alfred /usr/local/bin/alfred
ENTRYPOINT ["alfred"]

