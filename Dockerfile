ARG VERSION

ARG NOT_LATEST=${VERSION#latest}
ARG BASE=${NOT_LATEST:+version}
ARG BASE=${BASE:-latest}

FROM scratch AS latest
ARG TARGETARCH
ADD --unpack=true https://github.com/iiojib/bblog/releases/latest/download/bblog-linux-${TARGETARCH}.tar.gz /usr/local/bin/

FROM scratch AS version
ARG TARGETARCH
ARG VERSION
ADD --unpack=true https://github.com/iiojib/bblog/releases/download/v${VERSION#v}/bblog-linux-${TARGETARCH}.tar.gz /usr/local/bin/

FROM ${BASE}
ENTRYPOINT ["bblog"]
