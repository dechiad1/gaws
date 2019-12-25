package cmd

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/dechiad1/gaws/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ec2Cmd)
	ec2Cmd.AddCommand(ec2ListCmd)
}

var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "ec2 related actions",
	Long:  `ec2 related actions`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		sess = util.Auth()
		svc = ec2.New(sess)
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var ec2ListCmd = &cobra.Command{
	Use:   "list",
	Short: "list ec2s that are deployed",
	Long:  `list ec2s that are deployed in your account associated with the region specified. Display the external IP address (if any), the private IP address and the security groups that are attached`,
	Run: func(cmd *cobra.Command, args []string) {
		input := &ec2.DescribeInstancesInput{}
		result, err := svc.DescribeInstances(input)
		if err != nil {
			panic(err)
		}

		du := util.SetHeaders("Name (from tag)", "Private IP", "Security Groups", "Public DNS")
		for _, reservation := range result.Reservations {
			for _, instance := range reservation.Instances {
				//instanceId := *instance.InstanceId
				privateIp := *instance.PrivateIpAddress
				publicDns := *instance.PublicDnsName
				var sg_string strings.Builder
				r := rune(',')
				for i, sg := range instance.SecurityGroups {
					sg_string.WriteString(*sg.GroupId)
					if i < (len(instance.SecurityGroups) - 1) {
						sg_string.WriteRune(r)
					}
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
				du.AddRow(name, privateIp, sg_string.String(), publicDns)
			}
		}
		du.PrintDisplay()
	},
}
