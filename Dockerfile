FROM index.tenxcloud.com/docker_library/alpine
MAINTAINER admin@acale.ph

ADD build/bin/smartsc /usr/local/bin/smartsc

EXPOSE 5000
ENTRYPOINT ["/usr/local/bin/smartsc"]
