FROM docker.io/library/golang:1.18 as builder

WORKDIR /
ADD https://julialang-s3.julialang.org/bin/linux/x64/1.7/julia-1.7.2-linux-x86_64.tar.gz ./
RUN tar -zxvf julia-1.7.2-linux-x86_64.tar.gz && \
    mv julia-1.7.2 /usr/local/julia &&  \
    ln -s /usr/local/julia/bin/julia /usr/local/bin && \
    rm -rf julia-1.7.2-linux-x86_64.tar.gz

WORKDIR /gocode/julia
COPY go.mod ./
COPY *.go ./

WORKDIR /gocode/julia
COPY examples ./examples/

WORKDIR /gocode/julia/examples/matrix-inversion
RUN go mod tidy
RUN go build -o matrix-inversion ./

FROM docker.io/library/ubuntu:20.04 as base
WORKDIR /
ADD https://julialang-s3.julialang.org/bin/linux/x64/1.7/julia-1.7.2-linux-x86_64.tar.gz ./
RUN tar -zxvf julia-1.7.2-linux-x86_64.tar.gz && \
    mv julia-1.7.2 /usr/local/julia &&  \
    ln -s /usr/local/julia/bin/julia /usr/local/bin && \
    rm -rf julia-1.7.2-linux-x86_64.tar.gz

WORKDIR /go-julia
COPY --from=builder /gocode/julia/examples/matrix-inversion/matrix-inversion ./
