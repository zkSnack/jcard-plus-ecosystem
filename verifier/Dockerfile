# Build Go code
FROM golang:latest as verifier

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
COPY verifier/main.go ./verifier/main.go
COPY verifier/go.sum ./verifier/go.sum
COPY verifier/go.mod ./verifier/go.mod

WORKDIR /build/verifier
RUN go mod tidy
RUN go build .

RUN mv /build/verifier/verifier /home/app/verifier
COPY verifier/keys /home/app/keys
COPY verifier/static /home/app/static

# Remove all raw code and keep only built artifacts
RUN rm -rf /build

RUN chown -R $APP_USER:$APP_USER /home/app
USER app:app
WORKDIR /home/app

# Command to run
ENTRYPOINT ["/home/app/verifier"]

EXPOSE 9090