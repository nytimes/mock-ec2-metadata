FROM "golang:1.12"

COPY . /app

RUN cd /app && make deps && make docker

EXPOSE 8111/tcp

CMD ["/app/bin/mock-ec2-metadata"]
