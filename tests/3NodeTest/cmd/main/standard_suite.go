package main

import (
	"fmt"
	"strconv"
)

type StandardSuite struct {
	c *Config
}

func (t *StandardSuite) Name() string {
	return "Standard"
}

func (t *StandardSuite) RunTests() {
	// Test Keygroups
	logNodeAction(t.c.nodeA, "Creating keygroup testing")
	t.c.nodeA.CreateKeygroup("testing", true, 0, false)

	logNodeAction(t.c.nodeA, "Deleting keygroup testing")
	t.c.nodeA.DeleteKeygroup("testing", false)

	logNodeAction(t.c.nodeA, "Deleting nonexistent keygroup")
	t.c.nodeA.DeleteKeygroup("trololololo", true)

	logNodeAction(t.c.nodeA, "Creating Keygroup KG1")
	t.c.nodeA.CreateKeygroup("KG1", true, 0, false)

	// Test Get/Put of a single node
	logNodeAction(t.c.nodeA, "Putting KG1Item/KG1Value into KG1")
	t.c.nodeA.PutItem("KG1", "KG1Item", "KG1Value", false)

	logNodeAction(t.c.nodeA, "Getting the value in KG1")

	resp := t.c.nodeA.GetItem("KG1", "KG1Item", false)

	if resp != "KG1Value" {
		logNodeFailure(t.c.nodeA, "resp is \"KG1Value\"", resp)
	}

	logNodeAction(t.c.nodeA, "Getting a Value from a nonexistent keygroup")
	t.c.nodeA.GetItem("trololoool", "wut", true)

	logNodeAction(t.c.nodeA, "Putting a Value into a nonexistent keygroup")
	t.c.nodeA.PutItem("nonexistentkeygroup", "item", "data", true)

	logNodeAction(t.c.nodeA, "Putting new value KG1Item/KG1Value2 into KG1")
	t.c.nodeA.PutItem("KG1", "KG1Item", "KG1Value2", false)

	logNodeAction(t.c.nodeA, "Getting the value in KG1")
	resp = t.c.nodeA.GetItem("KG1", "KG1Item", false)
	if resp != "KG1Value2" {
		logNodeFailure(t.c.nodeA, "resp is \"KG1Value2\"", resp)
	}

	// test scanning
	logNodeAction(t.c.nodeA, "Creating kg scantest")
	t.c.nodeA.CreateKeygroup("scantest", true, 0, false)
	numItems := 20
	scanStart := 5
	scanRange := 10
	// 2. put in a bunch of items
	ids := make([]string, numItems)
	vals := make([]string, numItems)

	for i := 0; i < 20; i++ {
		vals[i] = "val" + strconv.Itoa(i)
		ids[i] = "id" + strconv.Itoa(i)
		t.c.nodeA.PutItem("scantest", ids[i], vals[i], false)
	}

	// 3. do a scan read
	// we expect [scanRange] amount of items, starting with [scanStart]
	startKey := "id" + strconv.Itoa(scanStart)

	items := t.c.nodeA.ScanItems("scantest", startKey, uint64(scanRange), false)

	expected := scanRange - scanStart

	if len(items) != expected {
		logNodeFailure(t.c.nodeA, fmt.Sprintf("%d items", expected), fmt.Sprintf("%d items", len(items)))
	}

	for i := 0; i < numItems; i++ {
		if i < scanStart || i >= scanStart+expected {
			continue
		}
		val, ok := items[ids[i]]
		if !ok {
			logNodeFailure(t.c.nodeA, fmt.Sprintf("%s is in returned items", ids[i]), fmt.Sprintf("%s is not in returned items", ids[i]))
			continue
		}
		if val != vals[i] {
			logNodeFailure(t.c.nodeA, fmt.Sprintf("item %s is %s", ids[i], vals[i]), fmt.Sprintf("item %s is %s", ids[i], items[ids[i]]))
			continue
		}
	}

	logNodeAction(t.c.nodeA, "Getting all Replicas that nodeA has")
	parsed := t.c.nodeA.GetAllReplica(false)
	// Example Response: map[string]string
	// {"nodeA": "1.2.3.4:5000", "nodeB": "4.5.6.7:4000"}
	// Test for nodeA

	if len(parsed) != 3 {
		logNodeFailure(t.c.nodeA, "GetAllReplica returns 3 nodes", fmt.Sprintf("%d", len(parsed)))
	}

	addr, ok := parsed[t.c.nodeA.ID]
	if !ok {
		logNodeFailure(t.c.nodeA, "GetAllReplica response contains nodeA", "nodeA not found")
	} else if addr != fmt.Sprintf("%s:%s", t.c.nodeAhost, t.c.nodeAhttpPort) {
		logNodeFailure(t.c.nodeA, "nodeA address is "+fmt.Sprintf("%s:%s", t.c.nodeAhost, t.c.nodeAhttpPort), addr)
	}

	addr, ok = parsed[t.c.nodeB.ID]
	if !ok {
		logNodeFailure(t.c.nodeA, "GetAllReplica response contains nodeB", "nodeB not found")
	} else if addr != fmt.Sprintf("%s:%s", t.c.nodeBhost, t.c.nodeBhttpPort) {
		logNodeFailure(t.c.nodeA, "nodeB address is "+fmt.Sprintf("%s:%s", t.c.nodeBhost, t.c.nodeBhttpPort), addr)
	}

	addr, ok = parsed[t.c.nodeC.ID]
	if !ok {
		logNodeFailure(t.c.nodeA, "GetAllReplica response contains nodeC", "nodeC not found")
	} else if addr != fmt.Sprintf("%s:%s", t.c.nodeChost, t.c.nodeChttpPort) {
		logNodeFailure(t.c.nodeA, "nodeC address is "+fmt.Sprintf("%s:%s", t.c.nodeChost, t.c.nodeChttpPort), addr)
	}
}

func NewStandardSuite(c *Config) *StandardSuite {
	return &StandardSuite{
		c: c,
	}
}
