FROM ubuntu:latest
USER root
RUN apt update && apt install -y apt-utils && apt upgrade -y && apt autoremove -y
RUN apt install -y ca-certificates net-tools openssl iputils-ping
RUN mkdir -p /app/log
COPY ./bin/speedtest_display_bin /app/speedtest_display_bin
RUN chmod +x /app/speedtest_display_bin
COPY ./static /app/static
COPY ./template.html /app/template.html
COPY ./entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh
WORKDIR /app
USER root
ENTRYPOINT ["bash", "./entrypoint.sh"]
EXPOSE 8080