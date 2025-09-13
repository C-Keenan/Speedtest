FROM ubuntu:latest
USER root
RUN apt update && apt install -y apt-utils && apt upgrade -y && apt autoremove -y
RUN apt install -y ca-certificates net-tools openssl iputils-ping
RUN mkdir -p /app/log
COPY ./bin/speedtest.5 /usr/local/bin/speedtest.5
COPY ./bin/speedtest.md /usr/local/bin/speedtest.md
COPY ./entrypoint.sh /app/entrypoint.sh
COPY ./speedtest.sh /app/speedtest.sh
RUN chmod +x /app/entrypoint.sh
RUN chmod +x /app/speedtest.sh
WORKDIR /app
USER root
ENTRYPOINT ["bash", "./entrypoint.sh"]