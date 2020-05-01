FROM debian:buster-slim

RUN set -ex; \
    apt-get update; \
    apt-get install --no-install-recommends -y \
        bash \
        ca-certificates \
        sudo \
        locales \
        procps \
        ; \
    apt-get autoclean; \
    rm -rf /var/lib/apt-lists/*;

SHELL [ "/bin/bash", "-Eexuo", "pipefail", "-c" ]

ENTRYPOINT [ "/bin/bash" ]
CMD [ "-c" "echo hello" ]
