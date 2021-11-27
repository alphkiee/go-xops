FROM openjdk:8-jdk-alpine

RUN rm -f /etc/localtime \
&& ln -sv /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone

ENV LANG en_US.UTF-8

ADD rocketmq-all-4.9.1-bin-release.tar.gz /usr/local/
RUN mv /usr/local/rocketmq-all-4.9.1-bin-release /usr/local/rocketmq-4.9.1 \
&& mkdir -p /data/rocketmq/store

CMD ["/bin/bash"]