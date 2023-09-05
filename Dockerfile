FROM golang:latest

WORKDIR /luqchain


RUN curl https://get.ignite.com/cli! | bash

COPY . /luqchain/

RUN ignite chain init && echo done

