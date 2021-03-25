FROM python

ENV PYTHONUNBUFFERED=1

RUN mkdir /app
WORKDIR /app
ADD . /app/
RUN pip install pika
EXPOSE 5000
CMD ["python", "/app/consumer.py"]