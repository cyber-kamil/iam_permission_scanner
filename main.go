package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func main() {
	// Create a new session with the AWS SDK
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a new IAM service client
	svc := iam.New(sess)

	// Open a file to write the output to
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error opening output file: ", err)
		os.Exit(1)
	}
	defer file.Close()

	// Get a list of all users
	result, err := svc.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		fmt.Println("Error getting list of users: ", err)
		os.Exit(1)
	}

	// Loop through each user and get their attached policies
	for _, user := range result.Users {
		fmt.Fprintf(file, "User: %s\n", *user.UserName)

		// Get a list of policies attached to the user
		attachedPolicies, err := svc.ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{
			UserName: user.UserName,
		})

		if err != nil {
			fmt.Fprintf(file, "\tError getting attached policies for user: %s - %v\n", *user.UserName, err)
			continue
		}

		// Print out each attached policy for the user
		for _, policy := range attachedPolicies.AttachedPolicies {
			fmt.Fprintf(file, "\tPolicy: %s\n", *policy.PolicyName)
		}
	}
}
