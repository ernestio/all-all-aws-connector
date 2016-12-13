/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	ecc "github.com/ernestio/ernest-config-client"
	"github.com/ernestio/ernestaws"
	"github.com/ernestio/ernestaws/elb"
	"github.com/ernestio/ernestaws/firewall"
	"github.com/ernestio/ernestaws/instance"
	"github.com/ernestio/ernestaws/nat"
	"github.com/ernestio/ernestaws/network"
	"github.com/ernestio/ernestaws/rdscluster"
	"github.com/ernestio/ernestaws/rdsinstance"
	"github.com/ernestio/ernestaws/route53"
	"github.com/ernestio/ernestaws/s3"
	"github.com/ernestio/ernestaws/vpc"
	"github.com/nats-io/nats"
)

var nc *nats.Conn
var natsErr error
var err error

func getEvent(m *nats.Msg) (n ernestaws.Event) {
	key := os.Getenv("ERNEST_CRYPTO_KEY")
	parts := strings.Split(m.Subject, ".")

	switch parts[0] {
	case "network":
		n = network.New(m.Subject, m.Data, key)
	case "nat":
		n = nat.New(m.Subject, m.Data, key)
	case "firewall":
		n = firewall.New(m.Subject, m.Data, key)
	case "vpc":
		n = vpc.New(m.Subject, m.Data, key)
	case "instance":
		n = instance.New(m.Subject, m.Data, key)
	case "elb":
		n = elb.New(m.Subject, m.Data, key)
	case "s3":
		n = s3.New(m.Subject, m.Data, key)
	case "route53":
		n = route53.New(m.Subject, m.Data, key)
	case "rds_cluster":
		n = rdscluster.New(m.Subject, m.Data, key)
	case "rds_instance":
		n = rdsinstance.New(m.Subject, m.Data, key)
	}

	return n
}

func eventHandler(m *nats.Msg) {
	var n ernestaws.Event
	if n = getEvent(m); n == nil {
		log.Println("Unrecognized event subject '" + m.Subject + "'")
		return
	}

	subject, data := ernestaws.Handle(&n)
	if err := nc.Publish(subject, data); err != nil {
		log.Println("Couldn't publish to nats")
	}
}

func main() {
	nc = ecc.NewConfig(os.Getenv("NATS_URI")).Nats()
	events := strings.Split(os.Getenv("CONNECTORS"), ",")
	if len(events) == 0 {
		log.Println("No connectors configured, please specify connectors on env var CONNECTORS")
		return
	}

	for _, subject := range events {
		fmt.Println("listening for " + subject)
		if _, err := nc.Subscribe(subject, eventHandler); err != nil {
			log.Println("Couldn't publish to nats")
		}
	}

	runtime.Goexit()
}
