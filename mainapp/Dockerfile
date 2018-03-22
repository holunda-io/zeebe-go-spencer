FROM golang:1.8

WORKDIR /app

ADD waitForZeebe.sh /app
ADD mainapp /app
ADD process /app/process
ADD zbctl /app
ADD config.toml /app

EXPOSE 8080

ENTRYPOINT ["./waitForZeebe.sh"]
