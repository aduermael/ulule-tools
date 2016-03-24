FROM golang:1.6

# Install redis
RUN cd /tmp; wget http://download.redis.io/releases/redis-3.0.5.tar.gz; tar xzf redis-3.0.5.tar.gz
RUN rm /tmp/redis-3.0.5.tar.gz
RUN mv /tmp/redis-3.0.5 /srv/redis-3.0.5
RUN cd /srv/redis-3.0.5; make

# Go package experimental vendoring
ADD ./go /go

ADD ./data /data
RUN chmod +x /data/start.sh
RUN ln -s /srv/redis-3.0.5/src/redis-server /bin/redis-server
RUN ln -s /srv/redis-3.0.5/src/redis-cli /bin/redis-cli

ENTRYPOINT ["/data/start.sh"]