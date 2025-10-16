FROM ubuntu:24.04 AS snap7-builder

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends \
        build-essential ca-certificates curl p7zip-full \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build
RUN curl -fL -o snap7.7z \
        "https://sourceforge.net/projects/snap7/files/1.4.2/snap7-full-1.4.2.7z/download" \
    && 7z x snap7.7z \
    && cd snap7-full-1.4.2/build/unix \
    && make -f x86_64_linux.mk all


FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive \
    PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    TZ=America/Sao_Paulo

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates curl gnupg tzdata \
        python3 python3-pip python3-venv \
        unixodbc \
    && curl -sSL https://packages.microsoft.com/keys/microsoft.asc \
        | gpg --dearmor -o /usr/share/keyrings/microsoft-prod.gpg \
    && echo "deb [arch=amd64,arm64,armhf signed-by=/usr/share/keyrings/microsoft-prod.gpg] https://packages.microsoft.com/ubuntu/24.04/prod noble main" \
        > /etc/apt/sources.list.d/mssql-release.list \
    && apt-get update \
    && ACCEPT_EULA=Y apt-get install -y --no-install-recommends msodbcsql17 msodbcsql18 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY --from=snap7-builder \
        /build/snap7-full-1.4.2/build/bin/x86_64-linux/libsnap7.so \
        /usr/lib/libsnap7.so
RUN ldconfig

WORKDIR /app

COPY requirements.txt ./
RUN pip install --break-system-packages -r requirements.txt

COPY config.py database.py plc_worker.py state.py stats_worker.py web.py main.py ./
COPY templates ./templates

EXPOSE 8080

RUN useradd --create-home --shell /usr/sbin/nologin balanca \
    && chown -R balanca:balanca /app
USER balanca

CMD ["python3", "-u", "main.py"]
