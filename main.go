package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
	"strconv"
)

func main() {
	// Read the queue name from the console
	queue := flag.String("q", "", "The name of the queue")
	// Read the queue name from the console
	topic := flag.String("t", "", "The name of the topic")
	topicArn := flag.String("a", "", "The ARN of the topic")
	// Parse the flags
	flag.Parse()

	if *queue == "" {
		fmt.Println("You must supply the name of a queue (-q QUEUE)")
		return
	}

	if *topic == "" {
		fmt.Println("You must supply the name of a topic (-t TOPIC)")
		return
	}

	if *topicArn == "" {
		fmt.Println("You must supply the ARN of a topic (-a TOPIC ARN)")
		return
	}

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	var accumulator Accumulator
	fmt.Println("Welcome to the accumulator!")
	fmt.Println("Write < exit > to exit")
	for {
		fmt.Println("Push a new number")

		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		if input == "exit" {
			break
		}

		// Convert input to number
		_, err = strconv.Atoi(input)
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		// Create an Amazon SQS service client
		snsClient := sns.NewFromConfig(cfg)

		// Publish message to the topic
		_, err = snsClient.Publish(
			context.TODO(),
			&sns.PublishInput{
				Message:  aws.String(input),
				TopicArn: aws.String(*topicArn),
			})

		if err != nil {
			fmt.Println("Error publishing to the topic:", err)
			return
		}

		// Create an Amazon SQS service client
		sqsClient := sqs.NewFromConfig(cfg)

		// Receive message
		message, err := sqsClient.ReceiveMessage(
			context.TODO(),
			&sqs.ReceiveMessageInput{
				QueueUrl: queue,
			})

		if err != nil {
			log.Fatal("Error happened while reading the message from the provided sqs")
		}

		msg := message.Messages[0]
		body := msg.Body
		// Variable to hold the unmarshalled struct
		var notification Notification

		// Unmarshal the JSON string into the struct
		err = json.Unmarshal([]byte(*body), &notification)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		// Convert string to int
		value, err := strconv.Atoi(notification.Message)
		if err != nil {
			// handle the error
			fmt.Println("Error converting string to int:", err)
			return
		}
		// Add to the accumulator
		accumulator.Add(value)

		// Delete the message
		_, err = sqsClient.DeleteMessage(
			context.TODO(),
			&sqs.DeleteMessageInput{
				QueueUrl:      queue,
				ReceiptHandle: msg.ReceiptHandle,
			})
		if err != nil {
			// handle the error
			fmt.Println("Error in deleting the notification", err)
			return
		}
		fmt.Println("Number now is:", accumulator.GetValue())
	}
}
