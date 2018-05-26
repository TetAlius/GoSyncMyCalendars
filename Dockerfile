FROM golang:1.10

WORKDIR src/github.com/TetAlius/GoSyncMyCalendars
COPY . .

RUN go get github.com/lib/pq
RUN go get github.com/google/uuid
RUN go install -v ./...

CMD ["GoSyncMyCalendars"]