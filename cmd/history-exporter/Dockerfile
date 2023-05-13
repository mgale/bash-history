# Using moderc sqlite module:
# https://pkg.go.dev/modernc.org/sqlite#section-readme
# This removes the C sqlite3 libaries from the image along with
# requiring CGO_ENABLED
#FROM gcr.io/distroless/static
From ubuntu:22.04
COPY output /
ENTRYPOINT [ "/history-exporter" ]
CMD ["stream", "-v"]


# # syntax = docker/dockerfile:1

# FROM golang:1.19 AS build

# RUN apt-get update && \
#     apt-get install -y build-essential git cmake \
#                        zlib1g-dev libevent-dev \
#                        libelf-dev llvm \
#                        clang libc6-dev-i386

# RUN mkdir /out
# RUN mkdir /src
# WORKDIR /src

# # Link asm/byteorder.h into eBPF
# RUN ln -s /usr/include/x86_64-linux-gnu/asm/ /usr/include/asm

# # Build libbpf as a static lib
# RUN git clone https://github.com/libbpf/libbpf-bootstrap.git && \
#     cd libbpf-bootstrap && \
#     git submodule update --init --recursive

# RUN cd libbpf-bootstrap/libbpf/src && \
#     make BUILD_STATIC_ONLY=y && \
#     make install BUILD_STATIC_ONLY=y LIBDIR=/usr/lib/x86_64-linux-gnu/

# ENV CGO_ENABLED=0
# COPY ../../go.* .
# RUN go mod download
# COPY . .

# RUN --mount=type=cache,target=/root/.cache/go-build \
# GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/execsnoop .

# FROM gcr.io/distroless/static AS bin-unix
# COPY --from=build /out/execsnoop /
# CMD [ "/execsnoop" ]
