FROM ubuntu
ENV GO_FILE='go1.16.linux-amd64.tar.gz'
ENV PATH=$PATH:/usr/local/go/bin
RUN apt-get update \
&& apt-get install -y git-core wget \
&& url=https://golang.org/dl/$GO_FILE \
&& wget $url \
&& tar -C /usr/local -xzf $GO_FILE \
&& rm -f $GO_FILE