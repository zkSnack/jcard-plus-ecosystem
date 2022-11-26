# Build Go code
FROM golang:latest as holder

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
COPY holder ./holder
COPY core-wallet ./core-wallet

WORKDIR /build/holder
RUN go mod tidy
RUN go build .

WORKDIR /build/holder/holder-front-end
RUN npm install
RUN npm run build


RUN mv /build/holder/holder /home/app/holder
RUN mv /build/holder/holder-front-end/build /home/app/holder-ui

COPY static/compiled-circuits /home/app/compiled-circuits
COPY static/js /home/app/js
COPY holder/config.yaml /home/app/config.yaml
RUN mkdir -p /home/app/tmp

# Remove all raw code and keep only built artifacts
RUN rm -rf /build

RUN chown -R $APP_USER:$APP_USER /home/app
USER app:app
WORKDIR /home/app

ENV CIRCUITS_PATH=/home/app/compiled-circuits/
ENV CIRCUITS_JS=/home/app/js/
ENV GIN_MODE=release
ENV UI_STATIC_DIR=./holder-ui/static
ENV UI_HTML_DIR=./holder-ui/index.html

# Command to run
ENTRYPOINT ["/home/app/holder"]

EXPOSE 8080