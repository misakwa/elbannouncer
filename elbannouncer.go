package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/elb"
)

type Announcer struct {
	InstanceId string // "instance-id"
	RegionId   string // "placement/availability-zone"
	elb        *elb.ELB
}

// TODO: Make requests concurrently?
func (sk *Announcer) LinkELB(elbNames ...string) (err error) {
	for _, name := range elbNames {
		_, err = sk.elb.RegisterInstancesWithLoadBalancer([]string{sk.InstanceId}, name)
	}
	return err
}

// TODO: Make requests concurrently?
func (sk *Announcer) UnlinkELB(elbNames ...string) (err error) {
	for _, name := range elbNames {
		_, err = sk.elb.DeregisterInstancesFromLoadBalancer([]string{sk.InstanceId}, name)
	}
	return err
}

func NewAnnouncer() (sk *Announcer, err error) {
	sk = &Announcer{}
	b, err := aws.GetMetaData("placement/availability-zone")
	if err != nil {
		return nil, err
	}
	sk.RegionId = string(b[:len(b)-1])
	b, err = aws.GetMetaData("instance-id")
	if err != nil {
		return nil, err
	}
	sk.InstanceId = string(b)

	if sk.RegionId == "" {
		sk.RegionId = aws.USEast.Name
	}
	auth, err := aws.GetAuth("", "", "", time.Time{})
	if err != nil {
		return nil, err
	}
	sk.elb = elb.New(auth, aws.Regions[sk.RegionId])
	return
}

func main() {
	help := fmt.Sprintf("%s link|unlink elb1 elb2 ...", os.Args[0])
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatalf("Too few arguments: %s", help)
	}
	cmd := args[0]
	elbs := args[1:]
	announcer, err := NewAnnouncer()
	if err != nil {
		log.Fatal(err)
	}
	switch cmd {
	case "link":
		err = announcer.LinkELB(elbs...)
		if err != nil {
			log.Fatal(err)
		}
	case "unlink":
		err = announcer.UnlinkELB(elbs...)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Command not recognized: %s, %s", cmd, help)
	}
}
