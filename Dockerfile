FROM ubuntu:latest
LABEL authors="sdeev"

ENTRYPOINT ["top", "-b"]
