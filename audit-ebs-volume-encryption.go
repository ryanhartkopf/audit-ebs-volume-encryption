package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"runtime/pprof"
	"sort"
)

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	regionPtr := flag.String("region", "", "the AWS region")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var sess *session.Session
	var err error

	// Create session with provided region, or fall back to environment default
	if regionPtr != nil {
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(*regionPtr),
		})
	} else {
		sess, err = session.NewSession()
	}
	if err != nil {
		fmt.Println("Error creating session ", err)
		return
	}

	// Create ec2 service
	ec2Svc := ec2.New(sess)

	// Describe all unencrypted volumes
	volumes, err := ec2Svc.DescribeVolumes(&ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("encrypted"),
				Values: []*string{
					aws.String("false"),
				},
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Extract instance IDs from DescribeVolumes output
	var instanceIds []*string
	for i := 0; i < len(volumes.Volumes); i++ {
		if volumes.Volumes[i].Attachments != nil {
			instanceIds = append(instanceIds, aws.String(*volumes.Volumes[i].Attachments[0].InstanceId))
		}
	}

	// Get name tags for instances
	instances, err := ec2Svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Create map of pointers to Name tags by instance ID
	idMap := map[string]*string{}
	for _, v := range instances.Reservations {
		var nameTag *string
		for _, t := range v.Instances[0].Tags {
			if *t.Key == "Name" {
				nameTag = t.Value
			}
		}
		idMap[*v.Instances[0].InstanceId] = nameTag
	}

	// Loop over list of volumes
	var data [][]string
	for i := 0; i < len(volumes.Volumes); i++ {

		// Variables
		volumeId := *volumes.Volumes[i].VolumeId
		var instanceId, device, tagName string

		// If volume is attached, set variables
		if volumes.Volumes[i].Attachments != nil {
			instanceId, device = *volumes.Volumes[i].Attachments[0].InstanceId, *volumes.Volumes[i].Attachments[0].Device
		}

		// If attached instance has a Name tag, set variable
		if idMap[instanceId] != nil {
			tagName = *idMap[instanceId]
		}

		// Append data to slice
		data = append(data, []string{volumeId, instanceId, device, tagName})
	}

	// Sort the data by instance name ascending
	data = stringsAscending(data, 3)

	// Create table - thank you Aleku Konko
	// https://github.com/olekukonko/tablewriter
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"volume_id", "instance_id", "device", "instance_name"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}

// Sorting is hard, sorting a slice of slices is harder
// so I stole the code below from Gyuho Lee @ AWS
// https://github.com/gyuho/learn/tree/master/doc/go_sort_algorithm#sort-table
var sortColumnIndex int

// sortByIndexAscending sorts two-dimensional strings in an ascending order, at a specified index.
type sortByIndexAscending [][]string

func (s sortByIndexAscending) Len() int {
	return len(s)
}

func (s sortByIndexAscending) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortByIndexAscending) Less(i, j int) bool {
	return s[i][sortColumnIndex] < s[j][sortColumnIndex]
}

// stringsAscending sorts two dimensional strings in an ascending order.
func stringsAscending(rows [][]string, idx int) [][]string {
	sortColumnIndex = idx
	sort.Sort(sortByIndexAscending(rows))
	return rows
}

// sortByIndexDescending sorts two-dimensional strings in an Descending order, at a specified index.
type sortByIndexDescending [][]string

func (s sortByIndexDescending) Len() int {
	return len(s)
}

func (s sortByIndexDescending) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortByIndexDescending) Less(i, j int) bool {
	return s[i][sortColumnIndex] > s[j][sortColumnIndex]
}

// stringsDescending sorts two dimensional strings in a descending order.
func stringsDescending(rows [][]string, idx int) [][]string {
	sortColumnIndex = idx
	sort.Sort(sortByIndexDescending(rows))
	return rows
}
