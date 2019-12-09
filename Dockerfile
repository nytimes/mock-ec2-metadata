FROM debian

WORKDIR /mock-ec2-metadata

COPY bin/mock-ec2-metadata_0.4.1_linux_amd64 /mock-ec2-metadata/mock-ec2-metadata 
COPY mock-ec2-metadata-config.json mock-ec2-metadata-config.json

ENTRYPOINT ["./mock-ec2-metadata"]