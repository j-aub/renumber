FROM python:3.12-alpine

LABEL org.opencontainers.image.source=https://github.com/j-aub/renumber
# LABEL org.opencontainers.image.description=""
LABEL org.opencontainers.image.licenses=GPL-3.0-or-later

WORKDIR /app

COPY ./app .

RUN python3 -m pip install --no-cache-dir -r requirements.txt
RUN rm requirements.txt

CMD ["gunicorn", "-t", "120", "-b", "0.0.0.0:80", "-w", "4", "app:app"]
