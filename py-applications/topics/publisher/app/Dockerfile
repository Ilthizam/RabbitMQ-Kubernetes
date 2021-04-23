# FROM python

# WORKDIR /app

# COPY requirements.txt requirements.txt
# RUN pip install -r requirements.txt


# EXPOSE 80
# CMD ["python", "/app/main.py"]

FROM python

ENV PYTHONUNBUFFERED=1

RUN mkdir /app
WORKDIR /app
ADD . /app/
RUN pip install -r requirements.txt


EXPOSE 5000
CMD ["python", "/app/publisher.py"]