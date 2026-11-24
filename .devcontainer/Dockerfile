FROM mcr.microsoft.com/devcontainers/go:1.27@sha256:74f45fc4c2abe0de9a9acdb4225ffb3e58e6b28f9af2e209afaa8920a3d3b623

RUN apt-get update \
    && apt-get install -y --no-install-recommends shellcheck \
    && rm -rf /var/lib/apt/lists/*

# renovate: datasource=custom.go-oldstable depName=go-oldstable
ARG GO_OLDSTABLE=1.26.11
ENV GO_OLDSTABLE=${GO_OLDSTABLE}
RUN mkdir "/usr/local/go${GO_OLDSTABLE}" \
    && curl -fsSL "https://go.dev/dl/go${GO_OLDSTABLE}.linux-$(dpkg --print-architecture).tar.gz" \
        | tar -xz -C "/usr/local/go${GO_OLDSTABLE}" --strip-components=1 \
    && ln -s "/usr/local/go${GO_OLDSTABLE}/bin/go" "/usr/local/bin/go${GO_OLDSTABLE}"
