FROM golang:1.8

WORKDIR /app

ADD waitForZeebe.sh /app
ADD mainapp/mainapp /app
ADD mainapp/process /app/process
ADD zbctl /app
ADD mainapp/config.toml /app

EXPOSE 8080

ENTRYPOINT ["./waitForZeebe.sh", "mainapp"]