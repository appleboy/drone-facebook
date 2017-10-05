FROM centurylink/ca-certs

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>"

ADD drone-facebook /

ENTRYPOINT ["/drone-facebook"]
