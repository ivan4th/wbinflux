FROM scratch
COPY wbinflux /
ENTRYPOINT ["/wbinflux"]
