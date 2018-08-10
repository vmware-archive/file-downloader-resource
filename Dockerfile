FROM alpine
RUN apk add --no-cache bash tzdata ca-certificates git openssh unzip zip gzip tar
RUN git config --global user.email "git@localhost"
RUN git config --global user.name "git"
COPY check-linux /opt/resource/check
COPY in-linux /opt/resource/in
COPY out-linux /opt/resource/out
