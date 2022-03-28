FROM docker.io/library/golang:1.18 as builder

WORKDIR /
COPY julia_download.sh ./
RUN ./julia_download.sh

WORKDIR /gocode/julia
COPY go.mod ./
COPY *.go ./

WORKDIR /gocode/julia
COPY examples ./examples/

WORKDIR /gocode/julia/examples/matrix-inversion
RUN go mod tidy
RUN go build -o matrix-inversion ./

WORKDIR /gocode/julia/examples/matrix-multiplication
RUN go mod tidy
RUN go build -o matrix-multiplication ./

WORKDIR /gocode/julia/examples/json-serialization
RUN go mod tidy
RUN go build -o json-serialization ./

FROM docker.io/library/ubuntu:20.04 as base
RUN apt-get update && apt-get install wget -y
WORKDIR /
COPY julia_download.sh ./
RUN ./julia_download.sh

# run julia code to install packages
COPY init.jl ./
RUN julia ./init.jl

WORKDIR /go-julia
COPY --from=builder /gocode/julia/examples/matrix-inversion/matrix-inversion ./
COPY --from=builder /gocode/julia/examples/matrix-multiplication/matrix-multiplication ./
COPY --from=builder /gocode/julia/examples/json-serialization/json-serialization ./
