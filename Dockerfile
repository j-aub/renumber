FROM rust:latest as builder

RUN rustup target add x86_64-unknown-linux-musl
RUN apt update && apt install -y musl-tools musl-dev

WORKDIR /app

COPY ./app ./
RUN cargo build --target x86_64-unknown-linux-musl --release

WORKDIR /gawk
ADD https://ftp.gnu.org/gnu/gawk/gawk-5.3.0.tar.xz ./gawk.tar.xz
RUN tar --strip-components=1 -xvf gawk.tar.xz

RUN ./configure LDFLAGS=-static
RUN make

# image
FROM scratch

COPY --from=builder /app/target/x86_64-unknown-linux-musl/release/app ./app
COPY ./app/templates ./templates/

COPY --from=builder /gawk/gawk ./gawk

LABEL org.opencontainers.image.source=https://github.com/j-aub/renumber
# LABEL org.opencontainers.image.description=""
LABEL org.opencontainers.image.licenses=GPL-3.0-or-later

ENV PATH=.
EXPOSE 80:8000/tcp
CMD ["/app"]
