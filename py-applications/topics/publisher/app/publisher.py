#!/usr/bin/env python
import pika
import sys
import flask
from flask import request




app = flask.Flask(__name__)
app.config["DEBUG"] = True
 

@app.route('/run', methods=['POST'])
def run():


    credentials = pika.PlainCredentials("guest", "guest")
    connection = pika.BlockingConnection(pika.ConnectionParameters("rabbitmq-0.rabbitmq.rabbits.svc.cluster.local", "5672", '/', credentials ))    
    channel = connection.channel()

    channel.exchange_declare(exchange='topic_logs', exchange_type='topic')

    routing_key = request.json['namespace']
    message = request.json['command']
    
    # data={
    #     "namespace":routing_key,
    #     "command":message
    # }

    channel.basic_publish(
        exchange='topic_logs', routing_key=routing_key, body=routing_key)

    print(" [x] Sent %r:%r" % (routing_key, message))

    connection.close()

    return "<h1>RabitMQ working</p>"

    
app.run()
