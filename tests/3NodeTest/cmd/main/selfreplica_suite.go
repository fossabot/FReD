package main

import (
	"time"
)

type SelfReplicaSuite struct {
	c *Config
}

func (t *SelfReplicaSuite) Name() string {
	return "Self Replication"
}

func (t *SelfReplicaSuite) RunTests() {
	// testing adding a node as a replica for a keygroup on itself
	logNodeAction(t.c.nodeB, "Create and populate a new keygroup to test pulling")
	t.c.nodeB.CreateKeygroup("pulltest", true, 0, false)
	t.c.nodeB.PutItem("pulltest", "item1", "val1", false)
	t.c.nodeB.PutItem("pulltest", "item2", "val2", false)

	logNodeAction(t.c.nodeA, "add nodeA as a replica to that keygroup and see if it pulls the needed data on its own (sleep 3s)")
	t.c.nodeA.AddKeygroupReplica("pulltest", t.c.nodeA.ID, 0, false)
	time.Sleep(3 * time.Second)
	// check if the items exist
	if res := t.c.nodeA.GetItem("pulltest", "item1", false); res != "val1" {
		logNodeFailure(t.c.nodeA, "val1", res)
	}
	if res := t.c.nodeA.GetItem("pulltest", "item2", false); res != "val2" {
		logNodeFailure(t.c.nodeA, "val2", res)
	}

	logNodeAction(t.c.nodeA, "Add an item on nodeA, check wheter it populates to nodeB")
	t.c.nodeA.PutItem("pulltest", "item3", "val3", false)
	// check if nodeB also gets that item
	if res := t.c.nodeB.GetItem("pulltest", "item3", false); res != "val3" {
		logNodeFailure(t.c.nodeB, "val3", res)
	}
}

func NewSelfReplicaSuite(c *Config) *SelfReplicaSuite {
	return &SelfReplicaSuite{
		c: c,
	}
}
