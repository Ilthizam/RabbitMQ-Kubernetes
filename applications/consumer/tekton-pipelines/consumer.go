package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

)

var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

var m = make(map[string]amqp.Delivery)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	consume(clientset)

}

func consume(clientset *kubernetes.Clientset) {

	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/")

	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	q, err := ch.QueueDeclare(
		"publisher", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	fmt.Println("Channel and Queue established")

	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to register consumer", err)
	}

	forever := make(chan bool)
	// jobList, _ := clientset.BatchV1().Jobs("rabbits").List(context.TODO(), metav1.ListOptions{})

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			// createDeployment(clientset)
			
			createJob(d)

			// jobList.Reset()
			// jobList, _ = clientset.BatchV1().Jobs("rabbits").List(context.TODO(), metav1.ListOptions{})

			// printJobs(clientset)
			// printNodes(clientset)
			// printNamespaces(clientset)
			// printJobsInDefaultNamespace(clientset)
			// printDeploymentsInDefaultNamespace(clientset)

			// for {
			// 	for _, job := range jobList.Items {

			// if job.Status.Active > 0 {
			// 	fmt.Printf("Running - %s\n", job.Name)
			// } else {
			// 	if job.Status.Succeeded > 0 {

			// 		fmt.Println("Deletion Started")
			// 		deleteJob(clientset, job, d)
			// 		fmt.Println("Deletion Ended")
			// 		jobList.Reset()
			// 		jobList, _ = clientset.BatchV1().Jobs("rabbits").List(context.TODO(), metav1.ListOptions{})

			// 	} else {
			// 		fmt.Printf("Failed - %s\n", job.Name)
			// 	}
			// }

			// 	}
			// 	jobList.Reset()
			// 	jobList, _ = clientset.BatchV1().Jobs("rabbits").List(context.TODO(), metav1.ListOptions{})
			// }

		}
	}()

	fmt.Println("Running...")
	<-forever

}

func spawnJobWatcher(clientset *kubernetes.Clientset) (<-chan (watch.Event), error) {
	watcher, err := clientset.BatchV1().Jobs("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return watcher.ResultChan(), nil
}

func deleteJob(clientset *kubernetes.Clientset, job batchv1.Job, d amqp.Delivery) {
	err := clientset.BatchV1().Jobs("rabbits").Delete(context.TODO(), job.Name, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	d.Ack(false)
	// fmt.Println("Jobs:")

	// for _, job := range jobList.Items {
	// 	fmt.Printf("- %d\n", job.Status.Succeeded)
	// }
	fmt.Printf("- %s", job.Name)
	fmt.Println("Deleted Job.")
	fmt.Println("")
}
func printJobs(clientset *kubernetes.Clientset) {
	jobList, _ := clientset.BatchV1().Jobs("rabbits").List(context.TODO(), metav1.ListOptions{})

	fmt.Println("Jobs:")

	for _, job := range jobList.Items {
		// fmt.Printf("- %d\n", job.Status.Succeeded)
		if job.Status.Succeeded == 1 {

			deleteJob(clientset, job, m[job.Name])
		}
	}

	fmt.Println("")
}

func printPods(clientset *kubernetes.Clientset) {
	podList, _ := clientset.CoreV1().Pods("rabbits").List(context.TODO(), metav1.ListOptions{})
	fmt.Println("Pods in rabbits:")

	for _, pod := range podList.Items {

		fmt.Printf("- %s\n", pod.Name)

	}

	fmt.Println("")
}

func printNodes(clientset *kubernetes.Clientset) {
	nodeList, _ := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	fmt.Println("Nodes:")

	for _, node := range nodeList.Items {
		fmt.Printf("- %s\n", node.Name)
	}

	fmt.Println("")
}

func printNamespaces(clientset *kubernetes.Clientset) {
	namespaceList, _ := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	fmt.Println("Namespaces:")

	for _, namespace := range namespaceList.Items {
		fmt.Printf("- %s\n", namespace.Name)
	}

	fmt.Println("")
}

func printJobsInDefaultNamespace(clientset *kubernetes.Clientset) {
	jobsClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)
	jobList, _ := jobsClient.List(context.TODO(), metav1.ListOptions{})

	if len(jobList.Items) > 0 {
		fmt.Println("Jobs in default namespace:")

		for _, job := range jobList.Items {
			fmt.Printf("- %s\n", job.Name)
		}
	} else {
		fmt.Println("No jobs in default namespace.")
	}

	fmt.Println("")
}

func printDeploymentsInDefaultNamespace(clientset *kubernetes.Clientset) {
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deploymentList, _ := deploymentsClient.List(context.TODO(), metav1.ListOptions{})

	if len(deploymentList.Items) > 0 {
		fmt.Println("Deployments in default namespace:")

		for _, deployment := range deploymentList.Items {
			fmt.Printf("- %s\n", deployment.Name)
		}
	} else {
		fmt.Println("No deployments in default namespace.")
	}

	fmt.Println("")
}

func createDeployment(clientset *kubernetes.Clientset) {
	deploymentsClient := clientset.AppsV1().Deployments("rabbits")
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

func createJob(d amqp.Delivery) {
	log.Printf("Create job funciton2")
	out,err := exec.Command("tkn" ,"task" ,"start" ,"echo-hello-world" ,"--showlog").Output() 


    if err != nil {
        log.Fatal(err)
    }
	log.Printf("%s",out)
	
}

func int32Ptr(i int32) *int32 { return &i }
