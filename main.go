package main

import (
        "fmt"
        "github.com/aws/aws-sdk-go/service/ec2"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/aws"
        "os"
        "strings"

)

var OWNER       = os.Getenv("OWNER")
var ACCOUNT_IDS = strings.Split(os.Getenv("ACCOUNT_IDS"), ",")
var AWS_REGION  = os.Getenv("AWS_DEFAULT_REGION")

func main() {
        if len(OWNER) == 0 {
                fmt.Println("Please specify owner account ID to accept peering connection")
                return
        }

        if len(ACCOUNT_IDS) <= 0 {
                fmt.Println("Please specify ACCOUNT IDs you own")
                return
        }

        if len(AWS_REGION) <= 0 {
                fmt.Printf("%s", "Please specify DEFAULT_AWS_REGION")
                return
        }
        fmt.Printf("%s\n", "Making VPC Peering Requests")

        svc := ec2.New(session.New(), &aws.Config{ Region: aws.String(AWS_REGION) })

        params := &ec2.DescribeVpcPeeringConnectionsInput{}
        resp, err := svc.DescribeVpcPeeringConnections(params)
        if err != nil {
                fmt.Println(err.Error())
                return
        }
        parseResponse(resp, svc)
        fmt.Printf("%s\n", "No VPC Peering Pending Request")
}

func parseResponse(resp *ec2.DescribeVpcPeeringConnectionsOutput, svc *ec2.EC2) {

        for _, v := range resp.VpcPeeringConnections {
                if isValidAccount(v.RequesterVpcInfo.OwnerId) == true && isOwner(v.AccepterVpcInfo.OwnerId) == true {
                        err := acceptPeeringRequest(v, svc)
                        if err != nil {
                                fmt.Println(err.Error())
                        }
                }
        }
}

func isValidAccount(vpc *string) bool {
        for _, v := range ACCOUNT_IDS {
                if v == *vpc {
                        return true
                }
        }
        return false
}

func isOwner(vpc *string) bool {
         if OWNER == *vpc {
                 return true
         }
         return false
}

func acceptPeeringRequest(v *ec2.VpcPeeringConnection, svc *ec2.EC2) error {
        if *v.Status.Code == "pending-acceptance" {
                params := &ec2.AcceptVpcPeeringConnectionInput{
                        VpcPeeringConnectionId: aws.String(*v.VpcPeeringConnectionId),
                }
                resp, err := svc.AcceptVpcPeeringConnection(params)
                if err != nil {
                        return err
                }
                fmt.Printf("VPC peering accepted")
                fmt.Printf("%s\n", resp)
        }
        return nil
}
