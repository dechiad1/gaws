package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/dechiad1/gaws/util"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

func init() {
	sgListGroups.PersistentFlags().StringVarP(&GroupNameFilter, "filter", "f", "", "filter sg based on if the value of this flag is in the name of the group")
	sgAddLocalGroup.PersistentFlags().StringVarP(&GroupIdInput, "group", "g", "", "group id of sg to add rule to")
	sgCmd.AddCommand(sgRemoveLocalGroup)
	sgCmd.AddCommand(sgListGroups)
	sgCmd.AddCommand(sgAddLocalGroup)
	rootCmd.AddCommand(sgCmd)
}

var (

	//flags
	GroupNameFilter string
	GroupIdInput    string
	GawsRuleName    string = "IP by gaws - current location"
)

//commands
var sgCmd = &cobra.Command{
	Use:   "sg",
	Short: "sg related actions",
	Long:  "list sgs, add an ip to sg ingress, remove sg ingress.. perhaps",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		sess = util.Auth()
		svc = ec2.New(sess)
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var sgRemoveLocalGroup = &cobra.Command{
	Use:   "rm",
	Short: "wipe the gaws rules that have been added",
	Long:  "remove all the gaws rules that have been added via and are identified with the GawsRuleName",
	Run: func(cmd *cobra.Command, args []string) {
		input := &ec2.DescribeSecurityGroupsInput{
			GroupIds: aws.StringSlice(args),
		}
		result, err := svc.DescribeSecurityGroups(input)
		if err != nil {
			fmt.Println("cant get any of the groups")
		}

		for _, group := range result.SecurityGroups {
			//fmt.Println(GawsRuleName)
			//fmt.Println("name: ", *group.GroupName)
			for _, permission := range group.IpPermissions {
				for _, ip := range permission.IpRanges {
					if *ip.Description == GawsRuleName {
						input := &ec2.RevokeSecurityGroupIngressInput{
							GroupId:    group.GroupId,
							IpProtocol: permission.IpProtocol,
							CidrIp:     ip.CidrIp,
							FromPort:   aws.Int64(22),
							ToPort:     aws.Int64(22),
						}
						req, resp := svc.RevokeSecurityGroupIngressRequest(input)
						err := req.Send()
						if err != nil {
							fmt.Println(err.Error())
							os.Exit(1)
						} else {
							fmt.Println(resp)
						}
					}
				}
			}
		}
	},
}

var sgAddLocalGroup = &cobra.Command{
	Use:   "add",
	Short: "add ingress rule for personal /32 ip: gaws sg add",
	Long:  "add your /32 as an ingress rule to a security group",
	Run: func(cmd *cobra.Command, args []string) {
		if GroupIdInput == "" {
			fmt.Println("group id not set. please add the group id with the '-g' flag")
			os.Exit(1)
		}
		ip := util.GetPublicIp()
		cidr := ip + "/32"
		fmt.Println(ip)

		input := &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String(GroupIdInput),
			IpPermissions: []*ec2.IpPermission{
				{
					FromPort:   aws.Int64(22),
					IpProtocol: aws.String("tcp"),
					IpRanges: []*ec2.IpRange{
						{
							CidrIp:      aws.String(cidr),
							Description: aws.String("IP by gaws - current location"),
						},
					},
					ToPort: aws.Int64(22),
				},
			},
		}

		result, err := svc.AuthorizeSecurityGroupIngress(input)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println(result)
		return
	},
}

type gEC2 struct {
	Client ec2iface.EC2API
}
 /*
 	Make an API call to aws with a possible filter on the 'group-name' value
 	return the [un]filtered security groups
  */
func (c *gEC2) listSGs(args []string) *ec2.DescribeSecurityGroupsOutput {
	value := fmt.Sprintf("*%s*", GroupNameFilter)
	name := "group-name"
	filter := &ec2.Filter{
		Name:   &name,
		Values: []*string{&value},
	}

	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{filter},
	}

	result, err := c.Client.DescribeSecurityGroups(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "InvalidGroupId.Malformed":
				fmt.Println("malformed sg group id u carb")
			case "InvalidGroup.NotFound":
				fmt.Println("group filter doesnt have any groups u fool")
			}
		}
		fmt.Println("cant get ")
	}

	return result
}

/*
	Format the security groups into output on the command line
 */
func printListedSGs(sgOutput *ec2.DescribeSecurityGroupsOutput) {
	fmt.Println("Security groups:")
	for _, group := range sgOutput.SecurityGroups {
		du := util.SetHeaders("Group Id", "Source", "From Port", "To Port")
		name := *group.GroupName
		id := *group.GroupId
		fmt.Println()
		fmt.Println("*******", name, "********")
		ipPermission := group.IpPermissions
		//ipPermission - top level object that contains the useful data for each ip rule
		for _, permission := range ipPermission {
			ipProtocol := *permission.IpProtocol
			//cidr blocks are found in the IpRanges object
			ipRanges := permission.IpRanges
			if len(ipRanges) > 0 {
				for _, ip := range ipRanges {
					if ipProtocol == "-1" {
						du.AddRow(id, *ip.CidrIp, "All", "All")
					} else {
						du.AddRow(id, *ip.CidrIp, strconv.FormatInt(*permission.FromPort, 10), strconv.FormatInt(*permission.ToPort, 10))
					}
				}
			}
			//SGs rules are found in the UserIdGroupPairs object
			groupPair := permission.UserIdGroupPairs
			if len(groupPair) > 0 {
				for _, pair := range groupPair {
					//if ipProtocol is -1, then connection is available on all ports! else - specify port range
					if ipProtocol == "-1" {
						du.AddRow(id, *pair.GroupId, "All", "All")
					} else {
						du.AddRow(id, *pair.GroupId, strconv.FormatInt(*permission.FromPort, 10), strconv.FormatInt(*permission.ToPort, 10))
					}
				}
			}
		}
		du.PrintDisplay()
	}
}

var sgListGroups = &cobra.Command{
	Use:   "list",
	Short: "list sg groups: gwas sg list <sg a> <sg b> <etc>",
	Long:  "list sg groups associated with the region one has set as an env variable",
	Run: func(cmd *cobra.Command, args []string) {
		c := gEC2{Client: svc}
		result := c.listSGs(args)
		//DEBUG: fmt.Print(result)
		printListedSGs(result)
	},
}
