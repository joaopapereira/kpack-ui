FROM node:10 as node-builder
WORKDIR /usr/src/app
COPY ui/package*.json ./
RUN npm install
COPY ui/ .
RUN npm run build

FROM golang as builder
ARG cmd=kpack-ui
COPY . /kpack-ui
WORKDIR /kpack-ui
RUN go build -o $cmd "./cmd/kpack-ui"

FROM ubuntu:18.04
COPY --from=builder "/kpack-ui/kpack-ui" "/kpack-ui/cmd"
RUN mkdir "/kpack-ui/conf"
COPY conf/app.prod.conf "/kpack-ui/conf/app.conf"
RUN mkdir -p /kpack-ui/ui/dist
COPY --from=node-builder "/usr/src/app/dist/*" "/kpack-ui/ui/dist/"

ENTRYPOINT ["/kpack-ui/cmd"]