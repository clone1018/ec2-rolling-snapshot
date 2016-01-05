package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"sort"
	"time"
)

type SnapshotTask struct {
	Volume    string
	Snapshots int
	Region    string
	Svc       *ec2.EC2
}

func (task *SnapshotTask) CreateSvc() {
	task.Svc = ec2.New(session.New(), &aws.Config{Region: aws.String(task.Region)})
}

func (task *SnapshotTask) CreateSnapshot(taskName string) (success bool, err error) {
	params := &ec2.CreateSnapshotInput{
		VolumeId:    aws.String(task.Volume),
		Description: aws.String(taskName),
	}

	resp, err := task.Svc.CreateSnapshot(params)
	if err != nil {
		return false, err
	}

	log.Printf("Creating snapshot: %s (%s)", *resp.SnapshotId, resp.StartTime.Format(time.RFC1123))

	return true, nil
}

func (task *SnapshotTask) DeleteOldSnapshots(taskName string) (success bool, err error) {
	params := &ec2.DescribeSnapshotsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("description"),
				Values: []*string{
					aws.String(taskName),
				},
			},
			{
				Name: aws.String("volume-id"),
				Values: []*string{
					aws.String(task.Volume),
				},
			},
		},
	}

	resp, err := task.Svc.DescribeSnapshots(params)
	if err != nil {
		return false, err
	}

	snapshots := resp.Snapshots
	sort.Sort(ByStartTime(snapshots))

	i := 1
	for _, snapshot := range snapshots {
		if i > task.Snapshots {
			snapshotId := *snapshot.SnapshotId

			log.Printf("Deleting snapshot: %s (%s)", snapshotId, snapshot.StartTime.Format(time.RFC1123))
			task.deleteSnapshot(snapshotId)
		}

		i++
	}

	return true, nil
}

func (task *SnapshotTask) deleteSnapshot(snapshot string) error {
	params := &ec2.DeleteSnapshotInput{
		SnapshotId: aws.String(snapshot),
	}
	_, err := task.Svc.DeleteSnapshot(params)

	return err
}

type ByStartTime []*ec2.Snapshot

func (a ByStartTime) Len() int      { return len(a) }
func (a ByStartTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByStartTime) Less(i, j int) bool {
	return a[i].StartTime.After(*a[j].StartTime)
}
