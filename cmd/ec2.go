package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	
	"github.com/dechiad1/gaws/util"
	"github.com/spf13/cobra"
  "github.com/aws/aws-sdk-go/service/ec2"
)

func init() {
	rootCmd.AddCommand(ec2Cmd)
	ec2Cmd.AddCommand(ec2ListCmd)
}

var ec2Cmd = &cobra.Command{
	Use: "ec2",
	Short: "ec2 related actions",
	Long: `ec2 related actions`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var ec2ListCmd = &cobra.Command{
	Use: "list",
	Short: "list ec2s that are deployed",
	Long: `list ec2s that are deployed in your account associated with the region specified. Display the external IP address (if any), the private IP address and the security groups that are attached`,
	Run: func(cmd *cobra.Command, args []string) {
		sess := util.Auth()
		svc := ec2.New(sess)

		input  := &ec2.DescribeInstancesInput{}
		result, err := svc.DescribeInstances(input)
		if err != nil {
			panic(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 20, 8, 0, '\t', 0)
		fmt.Fprintf(w, "\n %s\t%s\t%s\t%s", "Name (from tag)", "Private IP", "Security Groups", "Public DNS")
		fmt.Fprintf(w, "\n %s\t%s\t%s\t%s", "---------------", "----------", "---------------", "----------")
		for _, reservation := range result.Reservations {
			for _, instance := range reservation.Instances {
				//instanceId := *instance.InstanceId
				privateIp := *instance.PrivateIpAddress
				publicDns := *instance.PublicDnsName
				securityGroups := make([]string, len(instance.SecurityGroups))
				for i, sg := range instance.SecurityGroups {
					securityGroups[i] = *sg.GroupId
				}
				var name string
				for _, tag := range instance.Tags {
					if *tag.Key == "Name" {
							name = *tag.Value
					}
				}
				if name == "" {
						name = "NONAME"
				}
				fmt.Fprintf(w, "\n %s\t%s\t%s\t%s", name, privateIp, securityGroups, publicDns)
			}
		}
		w.Flush()
	},
}

