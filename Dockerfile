FROM golang:1.21-bookworm as builder

WORKDIR /gawk
ADD https://ftp.gnu.org/gnu/gawk/gawk-5.3.0.tar.xz ./gawk.tar.xz
RUN apt update \
	&& apt -y install xz-utils \
	&& tar --strip-components=1 -xvf gawk.tar.xz \
	&& ./configure LDFLAGS=-static \
	&& make

WORKDIR /renumber
COPY . /renumber/
# CGO_ENABLED=0 builds renumber statically which we need to run in a
# scratch container
RUN CGO_ENABLED=0 go build -ldflags "-s -w"

# image
FROM scratch

COPY --from=builder /renumber/renumber /renumber
# go is picky
# https://pkg.go.dev/os/exec#hdr-Executables_in_the_current_directory
COPY --from=builder /gawk/gawk /usr/bin/gawk

LABEL org.opencontainers.image.source=https://github.com/j-aub/renumber
# LABEL org.opencontainers.image.description=""
LABEL org.opencontainers.image.licenses=GPL-3.0-or-later

EXPOSE 80:8000/tcp
CMD ["/renumber"]
