package main

import (
        "fmt"
        "github.com/aws/aws-sdk-go/service/ec2"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/aws"
        "github.com/urfave/cli"
        "os"
)

var OWNER string
var AWS_REGION string
var ACCOUNT_IDS cli.StringSlice

func main() {

        app := cli.NewApp()
        app.Name = "vpc-peering-connection"
        app.Version = "0.0.1"
        cli.VersionPrinter = func(c *cli.Context) {
                fmt.Fprintf(c.App.Writer, "version=%s\n", c.App.Version)
        }
        app.Flags = []cli.Flag {
                cli.StringFlag{
                Name  : "owner-id, o",
                Usage : "owner-id that accept VPC Peering Request connection",
                EnvVar: "aws_account_id",
                Destination: &OWNER,
                },
                cli.StringSliceFlag{
                Name  : "list-aws-account-ids, l",
                Usage : "List of AWS account IDs to accept VPC Peering Request connection from",
                Value : &ACCOUNT_IDS,
                },
                cli.StringFlag{
                Name  : "DEFAULT AWS REGIONS, r",
                Usage : "Default AWS Region",
                EnvVar: "DEFAULT_AWS_REGION",
                Destination: &AWS_REGION,
                },
        }

        app.Action = func(c *cli.Context) error {
                if len(os.Args) < 2 {
                        cli.ShowAppHelp(c)
                        os.Exit(1)
                        return nil
                }
                fmt.Printf("%s\n", "Making VPC Peering Requests")

                svc := ec2.New(session.New(), &aws.Config{ Region: aws.String(AWS_REGION) })

                params := &ec2.DescribeVpcPeeringConnectionsInput{}
                resp, err := svc.DescribeVpcPeeringConnections(params)
                if err != nil {
                        fmt.Println(err.Error())
                        return err
                }
                parseResponse(resp, svc)
                fmt.Printf("%s\n", "No VPC Peering Pending Request")
                return nil
        }
        app.Run(os.Args)
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
