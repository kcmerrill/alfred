FROM ubuntu
COPY bin/alfred /usr/local/bin/alfred
ENTRYPOINT ["alfred"]

