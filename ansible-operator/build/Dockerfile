FROM registry.access.redhat.com/ubi8/ubi-minimal:latest


ENV OPERATOR=/usr/local/bin/ansible-operator \
    USER_UID=1001 \
    USER_NAME=ansible-operator
ENV PYMOL_VERSION 1.8.2.0
# install operator binary
COPY build/_output/bin/ansible-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
