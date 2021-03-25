# #!/usr/bin/env python

import pika
import flask

app = flask.Flask(__name__)
app.config["DEBUG"] = True


@app.route('/<task_name>', methods=['GET'])
def home(task_name):
    
    credentials = pika.PlainCredentials("guest", "guest")
    # connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
    connection = pika.BlockingConnection(pika.ConnectionParameters("rabbitmq-0.rabbitmq.rabbits.svc.cluster.local", "5672", '/', credentials ))
    channel = connection.channel()
    channel.queue_declare(queue='hello')
    channel.basic_publish(exchange='',
                      routing_key='hello',
                      body=task_name)
                      
    print(" [x] Sent 'Hello World!'")

    connection.close()

    return "<h1>RabitMQ working</p>"

app.run()
