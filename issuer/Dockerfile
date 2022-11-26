# Build Go code
FROM golang:latest as issuer

RUN curl -sL https://deb.nodesource.com/setup_16.x -o /tmp/nodesource_setup.sh
RUN chmod +x /tmp/nodesource_setup.sh
RUN /tmp/nodesource_setup.sh
RUN apt update && apt install -y nodejs && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
ENV APP_USER=app
ENV APP_UID=1001
RUN adduser -u $APP_UID $APP_USER --disabled-password --gecos "ZKSnacks,RoomNumber,WorkPhone,HomePhone"
ENV NPM_CONFIG_PREFIX=/home/app/node/.npm-global
RUN npm install -g snarkjs@latest
ENV PATH=${PATH}:/home/app/node/.npm-global/bin

WORKDIR /build
COPY issuer ./issuer
COPY core-issuer ./core-issuer
COPY core-wallet ./core-wallet

WORKDIR /build/issuer
RUN go mod tidy
RUN go build .

RUN mv /build/issuer/issuer /home/app/issuer
COPY static/compiled-circuits /home/app/compiled-circuits
COPY static/js /home/app/js
COPY data /home/app/data
COPY issuer/config.yaml /home/app/config.yaml
RUN mkdir -p /home/app/tmp

# Remove all raw code and keep only built artifacts
RUN rm -rf /build

RUN chown -R $APP_USER:$APP_USER /home/app
USER app:app
WORKDIR /home/app

ENV CLAIM_SCHEMA_DIR=/home/app/claim-schemas/
ENV ISSUER_DATA_DIR=/home/app/data
ENV CIRCUITS_PATH=/home/app/compiled-circuits/
ENV CIRCUITS_JS=/home/app/js/
ENV GIN_MODE=release

# Command to run
ENTRYPOINT ["/home/app/issuer"]

EXPOSE 8090