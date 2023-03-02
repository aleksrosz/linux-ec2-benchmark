package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"strings"
)

var (
	instanceTypesArray []types.InstanceType
	cfg                aws.Config
	instanceID         string
	instanceIP         string
	privateKey         string
	instanceType       string
)

// TODO
func readPubKey(file string) ssh.AuthMethod {
	var key ssh.Signer
	var err error
	var b []byte
	b, err = os.ReadFile(file)
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

func connectSSH(privateKey string, instanceIP string) {
	config := &ssh.ClientConfig{
		User: "ec2-user",
		Auth: []ssh.AuthMethod{
			readPubKey(privateKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// connect to ssh server
	conn, err := ssh.Dial("tcp", instanceIP+":22", config)
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
		"sudo yum -y install git gcc make automake libtool openssl-devel ncurses-compat-libs",
		"wget http://repo.mysql.com/mysql-community-release-el7-5.noarch.rpm",
		"sudo rpm -ivh mysql-community-release-el7-5.noarch.rpm",
		"sudo yum -y update",
		"sudo yum -y install mysql-community-devel mysql-community-client mysql-community-common",
		"git clone https://github.com/akopytov/sysbench",
		"cd sysbench",
		"./autogen.sh",
		"./configure",
		"make",
		"sudo make install",
		"sysbench cpu --threads=1 --cpu-max-prime=2500000 run",
	}
	command := strings.Join(commands, "; ")

	var buff bytes.Buffer
	session.Stdout = &buff
	err = session.Run(command)
	if err != nil {
		log.Fatal(err)
	}

	logging := false
	var buffer2 []byte
	scanner := bufio.NewScanner(&buff)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "sysbench 1.1.0") == true {
			logging = true
		}
		if logging == true {
			//log.Printf(scanner.Text())
			buffer2 = append(buffer2, scanner.Text()...)
			buffer2 = append(buffer2, '\n')
			err := os.WriteFile("./sysbench"+string(instanceType)+".txt", buffer2, 0600)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func readResultsFromFiles() {
	database := New()

	// Files number in ./results directory
	files, err := os.ReadDir("./results")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(files); i++ {
		result := readFile("sysbenchm5.large.txt")
		if err != nil {
			log.Fatal(err)
		}
		database.Add(result)
		foobar, ok := database.Get(0)
		if ok {
			fmt.Println(foobar.instanceName)
			fmt.Println(foobar.cpuSpeed)
		}

	}

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
//
//	c is the context of the method call, which includes the AWS Region.
//	api is the interface that defines the method call.
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a RunInstancesOutput object containing the result of the service call and nil.
//	Otherwise, nil and an error from the call to RunInstances.
func MakeInstance(c context.Context, api EC2CreateInstanceAPI, input *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error) {
	return api.RunInstances(c, input)
}

// MakeTags creates tags for an Amazon Elastic Compute Cloud (Amazon EC2) instance.
// Inputs:
//
//	c is the context of the method call, which includes the AWS Region.
//	api is the interface that defines the method call.
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a CreateTagsOutput object containing the result of the service call and nil.
//	Otherwise, nil and an error from the call to CreateTags.
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

	//TODO flags
	/*
		name := flag.String("n", string(InstanceType), "The name of the tag to attach to the instance")
		value := flag.String("v", string(InstanceType), "The value of the tag to attach to the instance")
		flag.Parse()

		if *name == "" || *value == "" {
			fmt.Println("You must supply a name and value for the tag (-n NAME -v VALUE)")
			return
		}

	*/

	//tags
	name := "test"
	value := "test2"

	client := ec2.NewFromConfig(cfg)

	minInstances := int32(1)
	maxInstances := int32(1)

	subnnett := "xxx"
	input := &ec2.RunInstancesInput{
		ImageId:          aws.String(ImageId),
		InstanceType:     InstanceType,
		MinCount:         &minInstances,
		MaxCount:         &maxInstances,
		KeyName:          aws.String(KeyName),
		SubnetId:         &subnnett,
		SecurityGroupIds: []string{"xxx"},
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
				Key:   &name,
				Value: &value,
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
	instanceID = *result.Instances[0].InstanceId
}

func getEC2ip(instanceID string) {
	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		fmt.Println("Got an error getting the instance:")
		fmt.Println(err)
	}
	instanceIP = *result.Reservations[0].Instances[0].PublicIpAddress
	fmt.Println(instanceIP)
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

// TODO if can't do something. Then delete instance.
func main() {
	instanceTypesArray = append(instanceTypesArray, types.InstanceTypeT2Nano)
	instanceTypesArray = append(instanceTypesArray, types.InstanceTypeT3Micro)
	/*

		for i := 0; i < len(instanceTypesArray); i++ {
			loadConfig()
			fmt.Println("Test for instance type: " + instanceTypesArray[i])
			instanceType = string(instanceTypesArray[i])
			createEC2Instance("ami-0a1ee2fb28fe05df3", instanceTypesArray[i], "home-PC")
			// TODO wait for instance to be running based on "Status check"
			time.Sleep(120 * time.Second)
			getEC2ip(instanceID)
			connectSSH("./home-PC.pem", instanceIP)
			terminateEC2Instace(instanceID)
		}
	*/

	// create in memory database for storing sysbenchResult structs
	resultsStore := New()

	for i := 0; i < len(instanceTypesArray); i++ {
		sysbenchResult1 := readFile(string("sysbench_" + instanceTypesArray[i] + ".txt"))
		resultsStore.Add(sysbenchResult1)
		appendToCSVFile()
	}
	resultsStore.Get(0)
	resultsStore.Get(1)

}
