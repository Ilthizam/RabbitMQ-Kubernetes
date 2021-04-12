FROM python

ENV PYTHONUNBUFFERED=1

RUN wget https://github.com/tektoncd/cli/releases/download/v0.17.0/tkn_0.17.0_Linux_x86_64.tar.gz
RUN tar xvzf tkn_0.17.0_Linux_x86_64.tar.gz -C /usr/local/bin/ tkn

RUN mkdir /app
WORKDIR /app
ADD . /app/
RUN pip install pika
EXPOSE 5000
CMD ["python", "/app/consumer.py"]