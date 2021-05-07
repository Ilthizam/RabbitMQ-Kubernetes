#!/usr/bin/env python
import pika
import sys,os
import subprocess as sp
from kubernetes import client, config


def main():

    config.load_incluster_config()


    api = client.CustomObjectsApi()

    credentials = pika.PlainCredentials("guest", "guest")
    connection = pika.BlockingConnection(pika.ConnectionParameters("rabbitmq-0.rabbitmq.rabbits.svc.cluster.local", "5672", '/', credentials ))    
    channel = connection.channel()

    channel.exchange_declare(exchange='topic_logs', exchange_type='topic')

    result = channel.queue_declare('', exclusive=True)
    queue_name = result.method.queue
    
    # binding_keys = sys.argv[1:]
    namespaces=['rabbits','team-a','team-b']
    # if not binding_keys:
    #     sys.stderr.write("Usage: %s [binding_key]...\n" % sys.argv[0])
    #     sys.exit(1)

    for namespace in namespaces:
        channel.queue_bind(
            exchange='topic_logs', queue=queue_name, routing_key=namespace)

    print(' [*] Waiting for logs. To exit press CTRL+C')


    def callback(ch, method, properties, body):
        print(" [x] %r:%r" % (method.routing_key, body))

        # namespace = request.json['namespace']
        # revision = request.json['command']
        # print(body.decode())
        # print(revision)

        api.create_namespaced_custom_object(
        group="tekton.dev",
        version="v1beta1",
        namespace=body.decode(),
        plural="taskruns",
        body={
            "apiVersion": "tekton.dev/v1beta1",
            "kind": "TaskRun",
            "metadata": {
                "generateName": "echo-hello-world-taskrun-",
                "namespace":body.decode()
            },
            "spec": {
                "serviceAccountName": "rabbitmq",
                "taskRef": {
                    "name":"echo-hello-world"
                }   
            },
        },
        )

        print("Resource created")

        # output = sp.getoutput(str('tkn taskrun list | grep Succeeded'))
        # tasks = len(output.splitlines())
        # if int(tasks)<2:

        # os.system(str(body.decode()))

            # ch.basic_ack(delivery_tag = method.delivery_tag)   


    channel.basic_consume(
        queue=queue_name, on_message_callback=callback, auto_ack=True)

    channel.start_consuming()

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print('Interrupted')
        try:
            sys.exit(0)
        except SystemExit:
            os._exit(0)