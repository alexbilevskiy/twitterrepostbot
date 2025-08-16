FROM buildpack-deps:bookworm-scm as local_builder

WORKDIR /

ENV TZ=Europe/Moscow

COPY --from=golang:1.24 /usr/local/go/ /usr/local/go/

ENV PATH="/usr/local/go/bin:${PATH}"

FROM local_builder

WORKDIR /app
COPY go.mod .
COPY go.sum .

COPY . .

RUN go build -o app cmd/twitterrepostbot.go

CMD ["/app/app"]
