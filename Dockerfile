FROM golang:1.21 as build
RUN mkdir /gips
WORKDIR /gips
ADD . .
RUN make linux
RUN chmod a+x /gips/bin/gips

FROM cgr.dev/chainguard/alpine-base:latest
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
LABEL org.label-schema.vcs-url="https://github.com/darron/gips"
RUN apk add --update --no-cache \
  ca-certificates && \
  rm -vf /var/cache/apk/*
WORKDIR /
COPY --from=build /gips/bin/gips /
ENTRYPOINT ["./gips", "service"]