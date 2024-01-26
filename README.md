# Go Accumulator Server



This Go application integrates AWS SNS (Simple Notification Service) and SQS (Simple Queue Service). 
It allows users to input numbers via the console, which are then published to an SNS topic. This topic forwards the messages to an SQS queue. 
The server polls messages from the queue, accumulates their values, and prints the running total.

## Prerequisites

Before running this application, ensure you have:

- An AWS account with access to SNS and SQS services.
- AWS CLI configured with the necessary permissions to interact with SNS and SQS.
- Go installed on your system.

## Setup

1. **AWS Configuration**: Make sure your AWS CLI is configured with credentials that have permission to publish to SNS topics and read from SQS queues. You can configure the AWS CLI by running `aws configure`.

2. **Create SNS Topic and SQS Queue**: Use the AWS documentation to create an SNS topic and an SQS queue. Then, subscribe to the SQS queue to the SNS topic.

3. **Install Dependencies**: If your application has external dependencies, install them using Go's package manager by running `go mod tidy`.

## Running the Application

Pass three command-line arguments when running the application:

- `-q`: The name of the SQS queue.
- `-t`: The name of the SNS topic.
- `-a`: The ARN of the SNS topic.


![image](https://github.com/muathejamil/accumulator/assets/27643048/9f6321d7-6adc-4e18-bba0-1136cbed67a1)

### Example Command

```sh
go run main.go -q YourQueueName -t YourTopicName -a YourTopicARN



