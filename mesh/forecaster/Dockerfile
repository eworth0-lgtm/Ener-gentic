FROM python:3.11-slim

WORKDIR /app

RUN pip install --no-cache-dir \
    numpy \
    grpcio \
    grpcio-tools \
    flask \
    pyarrow \
    pyarrow-flight

COPY proto/ /app/proto/
COPY forecaster.py /app/

RUN python -m grpc_tools.protoc \
    -I/app/proto \
    --python_out=/app \
    --grpc_python_out=/app \
    /app/proto/core.proto

EXPOSE 50051 8815

CMD ["python", "forecaster.py"]
