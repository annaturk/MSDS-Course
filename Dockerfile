FROM golang:1.21.6
WORKDIR /app
RUN go mod init msds-course
COPY *.go coursedata.csv ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /msds-course

# FROM gcr.io/distroless/base-debian11
# WORKDIR /
# COPY --from=builder /msds-courses.go /coursehandlers.go
ENV PORT 8080
# USER nonroot:nonroot
CMD ["/msds-course"]
