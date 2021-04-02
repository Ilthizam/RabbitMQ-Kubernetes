  
#!/usr/bin/env python
import pika, sys, os
import subprocess as sp

def main():
    credentials = pika.PlainCredentials("guest", "guest")
    # connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
    connection = pika.BlockingConnection(pika.ConnectionParameters("rabbitmq-0.rabbitmq.rabbits.svc.cluster.local", "5672", '/', credentials ))    
    channel = connection.channel()

    channel.queue_declare(queue='hello')

    def callback(ch, method, properties, body):
        print(" [x] Received %r" % body.decode())


        output = sp.getoutput(str('tkn taskrun list | grep Succeeded'))
        tasks = len(output.splitlines())
        if int(tasks)<2:
            os.system(str(body.decode()))
            ch.basic_ack(delivery_tag = method.delivery_tag)   


        
    
    channel.basic_consume(queue='hello', on_message_callback=callback, auto_ack=False)

    print(' [*] Waiting for messages. To exit press CTRL+C')
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