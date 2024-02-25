FROM golang:1.21.6 as builder
WORKDIR /app
RUN go mod init msds-course
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /msds-course

FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=builder /msds-courses /msds-courses
ENV PORT 1234
USER nonroot:nonroot
CMD ["/msds-course"]
