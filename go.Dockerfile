FROM golang:alpine AS builder

WORKDIR /usr/src/app
COPY ./go .

RUN go build -o main .

# I can't figure out how to make this work on scratch image.
# With scratch image, the build fails with an infinite loop.
FROM alpine

WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/main .

EXPOSE 8080
CMD [ "/usr/src/app/main" ]