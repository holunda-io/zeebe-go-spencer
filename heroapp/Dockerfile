FROM golang:1.8

WORKDIR /app

ADD waitForZeebe.sh /app
ADD zbctl /app
ADD heroapp/heroapp /app

EXPOSE 8080

ENTRYPOINT ["./waitForZeebe.sh", "heroapp"]