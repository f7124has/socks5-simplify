FROM debian:latest
RUN mkdir -p /app
COPY client.elf /app/client.elf
COPY server.elf /app/server.elf
RUN chmod +x /app/*
CMD ["/app/client.elf"]
