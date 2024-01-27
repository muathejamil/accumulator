package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"log"
	"strconv"
)

func main() {
	// Read the queue name, topic name, and the topic arn from the console
	queue, topicArn, err := ReadAndValidateArguments()
	if err != nil {
		return
	}

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Clients
	snsClient := sns.NewFromConfig(cfg)
	sqsClient := sqs.NewFromConfig(cfg)

	snsPublisher := &AWSSNSPublisher{Client: snsClient}
	sqsReceiver := &AWSSQSReceiver{Client: sqsClient}

	var accumulator Accumulator
	fmt.Println("Welcome to the accumulator!")
	fmt.Println("Write < exit > to exit")
	for {
		fmt.Println("Push a new number")
		input, err := ReadInput()
		if err != nil {
			break
		}

		// Publish message to the topic
		err = Publish(snsPublisher, input, topicArn)
		if err != nil {
			return
		}

		// Receive message
		message, err := ReadMsgFromSQS(sqsReceiver, queue)
		if err != nil {
			return
		}

		msg := message.Messages[0]
		value, err := GetIntMessage(msg)
		if err != nil {
			return
		}

		// Add to the accumulator
		accumulator.Add(value)

		// Delete the message
		if DeleteMsg(sqsReceiver, queue, msg) != nil {
			return
		}
		fmt.Println("Number now is:", accumulator.GetValue())
	}
}

// ReadAndValidateArguments read and validate the required arguments
func ReadAndValidateArguments() (*string, *string, error) {
	queue := flag.String("q", "", "The name of the queue")
	topic := flag.String("t", "", "The name of the topic")
	topicArn := flag.String("a", "", "The ARN of the topic")
	// Parse the flags
	flag.Parse()

	if *queue == "" {
		fmt.Println("You must supply the name of a queue (-q QUEUE)")
		return nil, nil, errors.New("you must supply the name of a queue (-q QUEUE)")
	}

	if *topic == "" {
		fmt.Println("You must supply the name of a topic (-t TOPIC)")
		return nil, nil, errors.New("you must supply the name of a topic (-t TOPIC)")
	}

	if *topicArn == "" {
		fmt.Println("You must supply the ARN of a topic (-a TOPIC ARN)")
		return nil, nil, errors.New("you must supply the ARN of a topic (-a TOPIC ARN)")
	}
	return queue, topicArn, nil
}

// DeleteMsg deletes the message
func DeleteMsg(sqsReceiver SQSReceiver, queue *string, msg types.Message) error {
	_, err := sqsReceiver.DeleteMessage(
		context.TODO(),
		&sqs.DeleteMessageInput{
			QueueUrl:      queue,
			ReceiptHandle: msg.ReceiptHandle,
		})
	if err != nil {
		// handle the error
		fmt.Println("Error in deleting the notification", err)
		return err
	}
	return nil
}

// ReadMsgFromSQS reads message from SQS
func ReadMsgFromSQS(sqsReceiver SQSReceiver, queue *string) (*sqs.ReceiveMessageOutput, error) {
	message, err := sqsReceiver.ReceiveMessage(
		context.TODO(),
		&sqs.ReceiveMessageInput{
			QueueUrl: queue,
		})

	if err != nil {
		log.Fatal("Error happened while reading the message from the provided sqs")
		return nil, err
	}
	return message, nil
}

// Publish publishes message to SNS
func Publish(snsPublisher SNSPublisher, input string, topicArn *string) error {
	_, err := snsPublisher.Publish(
		context.TODO(),
		&sns.PublishInput{
			Message:  aws.String(input),
			TopicArn: aws.String(*topicArn),
		})

	if err != nil {
		fmt.Println("Error publishing to the topic:", err)
		return err
	}
	return nil
}

// GetIntMessage convert message to its integer value
func GetIntMessage(message types.Message) (int, error) {
	body := message.Body
	// Variable to hold the unmarshalled struct
	var notification Notification

	// Unmarshal the JSON string into the struct
	err := json.Unmarshal([]byte(*body), &notification)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return 0, err
	}
	// Convert string to int
	value, err := strconv.Atoi(notification.Message)
	if err != nil {
		// handle the error
		fmt.Println("Error converting string to int:", err)
		return 0, err
	}
	return value, err
}

// IsNumber validates if entered string is number
func IsNumber(input string) bool {
	_, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return false
	}
	return true
}

// ReadInput reads the input from console and validate
func ReadInput() (string, error) {
	var input string
	_, err := fmt.Scan(&input)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return "", err
	}

	if input == "exit" {
		return "", errors.New("exit")
	}

	// Validate if input is valid number
	if !IsNumber(input) {
		return "", errors.New("invalid input. Please enter a numeric value")
	}
	return input, nil
}
