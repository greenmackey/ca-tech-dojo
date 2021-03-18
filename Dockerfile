FROM golang:1.16.2-buster
RUN apt-get update\
&& apt-get install -y git-core
CMD bash