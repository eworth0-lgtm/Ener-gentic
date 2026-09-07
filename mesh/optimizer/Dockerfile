FROM rust:1.75-slim

WORKDIR /app

RUN apt-get update && apt-get install -y \
    protobuf-compiler \
    libprotobuf-dev \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

COPY proto/ /app/proto/
COPY Cargo.toml /app/
COPY src/ /app/src/
COPY cpp_sim/ /app/cpp_sim/

RUN cd /app/cpp_sim && make

RUN cargo build --release

EXPOSE 50052

CMD ["./target/release/optimizer"]
