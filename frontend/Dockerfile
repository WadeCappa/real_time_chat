# syntax=docker/dockerfile:1

FROM python:3.9.20

WORKDIR /app

RUN pip3 install flask
RUN pip3 install gunicorn

COPY . ./

EXPOSE 8000

CMD ["gunicorn", "-w", "1", "-b", "0.0.0.0", "hello_world:app"]
