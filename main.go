package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws/client"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"strings"
)

var cfg aws.Config

type EC2 struct {
	*client.Client
}

func readPubKey(file string) ssh.AuthMethod {
	var key ssh.Signer
	var err error
	var b []byte
	b, err = ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err) //"failed to read public key")
	}
	if !strings.Contains(string(b), "ENCRYPTED") {
		key, err = ssh.ParsePrivateKey(b)
		if err != nil {
			log.Fatal(err) //"failed to parse public key")
		}
	} else {
		keyPass := "password"
		// Decrypt the key with the passphrase
		key, err = ssh.ParsePrivateKeyWithPassphrase(b, []byte(keyPass))
		if err != nil {
			log.Fatal(err) // "failed to parse password-protected private key"
		}
	}
	return ssh.PublicKeys(key)
}

func connectSSH() {
	config := &ssh.ClientConfig{
		User: "xxx",
		Auth: []ssh.AuthMethod{
			readPubKey("./xxx.pem"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// connect to ssh server
	conn, err := ssh.Dial("tcp", "xxx:22", config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	/*
		// configure terminal mode
		modes := ssh.TerminalModes{
			ssh.ECHO: 0, // supress echo

		}
		// run terminal session
		if err := session.RequestPty("xterm", 50, 80, modes); err != nil {
			log.Fatal(err)
		}

		// start remote shell
		if err := session.Shell(); err != nil {
			log.Fatal(err)
		}

	*/

	// Create a single command that is semicolon seperated
	commands := []string{
		"sysbench cpu --threads=1 --cpu-max-prime=50000 run",
	}
	command := strings.Join(commands, "; ")

	var buff bytes.Buffer
	session.Stdout = &buff
	err = session.Run(command)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buff.String())
}

// EC2CreateInstanceAPI defines the interface for the RunInstances and CreateTags functions.
// We use this interface to test the functions using a mocked service.
type EC2CreateInstanceAPI interface {
	RunInstances(ctx context.Context,
		params *ec2.RunInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)

	CreateTags(ctx context.Context,
		params *ec2.CreateTagsInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error)
}

// MakeInstance creates an Amazon Elastic Compute Cloud (Amazon EC2) instance.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a RunInstancesOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to RunInstances.
func MakeInstance(c context.Context, api EC2CreateInstanceAPI, input *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error) {
	return api.RunInstances(c, input)
}

// MakeTags creates tags for an Amazon Elastic Compute Cloud (Amazon EC2) instance.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a CreateTagsOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to CreateTags.
func MakeTags(c context.Context, api EC2CreateInstanceAPI, input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return api.CreateTags(c, input)
}

func terminateEC2Instace(instanceID string) {
	client := ec2.NewFromConfig(cfg)

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceID},
	}

	_, err := client.TerminateInstances(context.TODO(), input)
	if err != nil {
		fmt.Println("Got an error terminating the instance:")
		fmt.Println(err)
		return
	}

	fmt.Println("Terminated instance with ID: " + instanceID)
}

func createEC2Instance(ImageId string, InstanceType types.InstanceType, KeyName string) {
	name := flag.String("n", "test", "The name of the tag to attach to the instance")
	value := flag.String("v", "test2", "The value of the tag to attach to the instance")
	flag.Parse()

	if *name == "" || *value == "" {
		fmt.Println("You must supply a name and value for the tag (-n NAME -v VALUE)")
		return
	}

	client := ec2.NewFromConfig(cfg)

	minInstances := int32(1)
	maxInstances := int32(1)

	input := &ec2.RunInstancesInput{
		ImageId:      aws.String(ImageId),
		InstanceType: InstanceType,
		MinCount:     &minInstances,
		MaxCount:     &maxInstances,
		KeyName:      aws.String(KeyName),
	}

	result, err := MakeInstance(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error creating an instance:")
		fmt.Println(err)
		return
	}

	tagInput := &ec2.CreateTagsInput{
		Resources: []string{*result.Instances[0].InstanceId},
		Tags: []types.Tag{
			{
				Key:   name,
				Value: value,
			},
		},
	}

	_, err = MakeTags(context.TODO(), client, tagInput)
	if err != nil {
		fmt.Println("Got an error tagging the instance:")
		fmt.Println(err)
		return
	}

	fmt.Println("Created tagged instance with ID " + *result.Instances[0].InstanceId)
}

func loadConfig() {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	var err error
	cfg, err = config.LoadDefaultConfig(context.TODO(),
		config.WithSharedCredentialsFiles(
			[]string{"./credentials", "data/credentials"},
		),
		config.WithSharedConfigFiles(
			[]string{"./config", "data/config"},
		),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
}

func main() {
	loadConfig()
	//createEC2Instance("xxx", types.InstanceTypeT2Micro, "xxx")
	//terminateEC2Instace("xxx")
	//connectSSH()
}
