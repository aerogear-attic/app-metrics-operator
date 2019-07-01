FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/app-metrics-operator \
    USER_UID=1001 \
    USER_NAME=app-metrics-operator

# install operator binary
COPY build/_output/bin/app-metrics-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
