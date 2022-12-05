FROM uhub.service.ucloud.cn/entropypool/alpine:3.16

WORKDIR /app

COPY chain-gateway /app

ENTRYPOINT ["/app/chain-gateway"]

EXPOSE 50580 50581 50582 50583

CMD ["run"]
