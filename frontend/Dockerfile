
FROM python:3.9.20

WORKDIR /usr/src/app

RUN pip3 install flask
RUN pip3 install gunicorn

COPY . ./

CMD ["gunicorn", "-w", "1", "hello_world:app"]
