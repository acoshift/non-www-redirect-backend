FROM acoshift/go-scratch

USER 65534:65534
COPY server /
ENTRYPOINT ["/server"]
