package badgerdb

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/go-errors/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const badgerDBPath = "./test.db"

var db *Storage

// TODO: better tests, maybe even for all packages that implement the Store interface?

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	fInfo, err := os.Stat(badgerDBPath)

	if err == nil {
		if !fInfo.IsDir() {
			panic(errors.Errorf("%s is not a directory!", badgerDBPath))
		}

		err = os.RemoveAll(badgerDBPath)
		if err != nil {
			panic(err)
		}
	}

	db = New(badgerDBPath)

	stat := m.Run()

	fInfo, err = os.Stat(badgerDBPath)

	if err == nil {
		if !fInfo.IsDir() {
			panic(errors.Errorf("%s is not a directory!", badgerDBPath))
		}

		err = os.RemoveAll(badgerDBPath)
		if err != nil {
			panic(err)
		}
	}

	os.Exit(stat)
}

func TestKeygroups(t *testing.T) {
	kg := "test-kg"
	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	exists := db.ExistsKeygroup(kg)
	if !exists {
		t.Fatal("Keygroup does not exist after creation")
	}

	err = db.DeleteKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	exists = db.ExistsKeygroup(kg)
	if exists {
		t.Fatal("Keygroup does still exist after deletion")
	}

}

func TestReadSome(t *testing.T) {
	kg := "test-kg-scan"
	updates := 10
	scanStart := 3
	scanRange := 5

	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	// 2. put in a bunch of items
	ids := make([]string, updates)
	vals := make([]string, updates)

	for i := 0; i < updates; i++ {
		ids[i] = "id" + strconv.Itoa(i)
		vals[i] = "val" + strconv.Itoa(i)

		err = db.Update(kg, ids[i], vals[i], false, 0)

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			t.Error(err)
		}

	}

	res, err := db.ReadSome(kg, "id"+strconv.Itoa(scanStart), uint64(scanRange))

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	assert.Len(t, res, scanRange)

	for i := scanStart; i < scanStart+scanRange; i++ {
		assert.Contains(t, res, ids[i])
		assert.Equal(t, res[ids[i]], vals[i])
	}
}

func TestReadAll(t *testing.T) {
	kg := "test-read-all"
	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg, "id-1", "data-1", false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg, "id-2", "data-2", false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg, "id-3", "data-3", false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	kg2 := "test-read-all-2"

	err = db.CreateKeygroup(kg2)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg2, "id-1", "data-1", false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg2, "id-2", "data-2", false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg2, "id-3", "data-3", false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	res, err := db.ReadAll(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	assert.Equal(t, "data-1", res["id-1"])
	assert.Equal(t, "data-2", res["id-2"])
	assert.Equal(t, "data-3", res["id-3"])

}

func TestIDs(t *testing.T) {
	kg := "test-ids"
	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg, "id-1", "data-1", false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg, "id-2", "data-2", false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg, "id-3", "data-3", false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	kg2 := "test-read-all-2"

	err = db.CreateKeygroup(kg2)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg2, "id-1", "data-1", false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg2, "id-2", "data-2", false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg2, "id-3", "data-3", false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	res, err := db.IDs(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	assert.Equal(t, []string{"id-1", "id-2", "id-3"}, res)

}

func TestItemExists(t *testing.T) {
	kg := "test-kg-item"
	id := "name"
	id2 := "name2"
	value := "value"

	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg, id, value, false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	ex := db.Exists(kg, id)
	if !ex {
		t.Error("exists says existing item doesn't exist")
	}

	ex = db.Exists(kg, id2)
	if ex {
		t.Error("exists says non-existent item exists")
	}

}

func TestItemGet(t *testing.T) {
	kg := "test-kg-item"
	id := "name"
	value := "value"

	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg, id, value, false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	retr, err := db.Read(kg, id)
	if err != nil {
		t.Error(err)
	}
	if retr != value {
		t.Errorf("Expected to get %s but got %s", value, retr)
	}
}

func TestItemDelete(t *testing.T) {
	kg := "test-kg-item-delete"
	id := "name"
	id2 := "name2"
	value := "value"

	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg, id, value, false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	retr, err := db.Read(kg, id)
	if err != nil {
		t.Error(err)
	}
	if retr != value {
		t.Errorf("Expected to get %s but got %s", value, retr)
	}

	err = db.Delete(kg, id)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	retr, err = db.Read(kg, id)
	if err == nil {
		t.Errorf("read a deleted item: %s", retr)
	}

	err = db.Delete(kg, id2)

	if err != nil {
		t.Error(err, "deleting non-existent keys should be allowed")
	}

}

func TestItemAfterDeleteKeygroup(t *testing.T) {
	kg := "test-kg-item-delete"
	id := "ndel"
	value := "vdel"

	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg, id, value, false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.DeleteKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	retr, err := db.Read(kg, id)
	if err == nil {
		t.Errorf("Expected an error, but got %s", retr)
	}
}

func TestExpiry(t *testing.T) {
	kg := "test-kg-item"
	id := "name"
	value := "value"

	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg, id, value, false, 10)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	retr, err := db.Read(kg, id)
	if err != nil {
		t.Error(err)
	}
	if retr != value {
		t.Errorf("Expected to get %s but got %s", value, retr)
	}

	time.Sleep(10 * time.Second)

	_, err = db.Read(kg, id)
	if err == nil {
		t.Error(err)
	}
}

func TestAppend(t *testing.T) {
	kg := "log"
	v1 := "value-1"
	v2 := "value-2"

	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	key1, err := db.Append(kg, v1, 0)

	if err != nil {
		t.Error(err)
	}
	if key1 != "0" {
		t.Errorf("Expected to get %s but got %s", "0", key1)
	}

	key2, err := db.Append(kg, v2, 0)

	if err != nil {
		t.Error(err)
	}
	if key2 != "1" {
		t.Errorf("Expected to get %s but got %s", "0", key2)
	}

	for i := 2; i < 100; i++ {
		v := "value-" + strconv.Itoa(i)
		key, err := db.Append(kg, v, 0)

		if err != nil {
			t.Error(err)
		}

		if key != strconv.Itoa(i) {
			t.Errorf("Expected to get %s but got %s", strconv.Itoa(i), key)
		}
	}
}

func TestConcurrentAppend(t *testing.T) {
	kg := "logconcurrent"
	concurrent := 4
	items := 100

	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	keys := make([]map[string]struct{}, concurrent)
	done := make(chan struct{})

	for i := 0; i < concurrent; i++ {
		keys[i] = make(map[string]struct{})
		go func(id int, keys *map[string]struct{}) {
			for j := 2; j < items; j++ {
				v := fmt.Sprintf("value-%d-%d", id, j)
				key, err := db.Append(kg, v, 0)

				if err != nil {
					t.Error(err.(*errors.Error).ErrorStack())
				}

				(*keys)[key] = struct{}{}
			}
			done <- struct{}{}
		}(i, &keys[i])
	}

	for i := 0; i < concurrent; i++ {
		<-done
	}

	for i, k := range keys {
		for key := range k {
			found := false
			for j := i + 1; j < concurrent; j++ {
				_, ok := keys[j][key]
				found = found || ok
			}
			assert.False(t, found, "key given out multiple times")
		}
	}

}

func TestTriggerNodes(t *testing.T) {
	kg := "kg1"

	err := db.CreateKeygroup(kg)

	if err != nil {
		t.Error(err)
	}

	t1 := "t1"
	t1host := "1.1.1.1:3000"

	t2 := "t2"
	t2host := "2.2.2.2:3000"

	t3 := "t3"
	t3host := "3.3.3.3:3000"

	err = db.AddKeygroupTrigger(kg, t1, t1host)

	if err != nil {
		t.Error(err)
	}

	err = db.AddKeygroupTrigger(kg, t1, t1host)

	if err != nil {
		t.Error(err)
	}

	err = db.AddKeygroupTrigger(kg, t2, t2host)

	if err != nil {
		t.Error(err)
	}

	tList, err := db.GetKeygroupTrigger(kg)

	if err != nil {
		t.Error(err)
	}

	log.Debug().Msgf("List of keygroup triggers: %#v", tList)

	if len(tList) != 2 {
		t.Error("not the right number of triggers for this keygroup")
	}

	if _, ok := tList[t1]; !ok {
		t.Error("t1 not in list of triggers for this keygroup")
	}

	if _, ok := tList[t2]; !ok {
		t.Error("t2 not in list of triggers for this keygroup")
	}

	if host := tList[t1]; host != t1host {
		t.Error("t1host not correct")
	}

	if host := tList[t2]; host != t2host {
		t.Error("t1host not correct")
	}

	err = db.DeleteKeygroupTrigger(kg, t1)

	if err != nil {
		t.Error(err)
	}

	tList, err = db.GetKeygroupTrigger(kg)

	if err != nil {
		t.Error(err)
	}

	log.Debug().Msgf("List of keygroup triggers: %#v", tList)

	if len(tList) != 1 {
		t.Error("not the right number of triggers for this keygroup")
	}

	if _, ok := tList[t2]; !ok {
		t.Error("t2 not in list of triggers for this keygroup")
	}

	if host := tList[t2]; host != t2host {
		t.Error("t1host not correct")
	}

	err = db.AddKeygroupTrigger(kg, t3, t3host)

	if err != nil {
		t.Error(err)
	}

	err = db.DeleteKeygroup(kg)

	if err != nil {
		t.Error(err)
	}

	tList, _ = db.GetKeygroupTrigger(kg)

	log.Debug().Msgf("List of keygroup triggers: %#v", tList)

	if len(tList) != 0 {
		t.Error("got keygroup triggers for nonexistent keygroup")
	}

}

func TestClose(t *testing.T) {
	kg := "test-kg-item"
	id := "name"
	value := "value"

	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = db.Update(kg, id, value, false, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	retr, err := db.Read(kg, id)
	if err != nil {
		t.Error(err)
	}
	if retr != value {
		t.Errorf("Expected to get %s but got %s", value, retr)
	}

	err = db.Close()

	if err != nil {
		t.Error(err)
	}

	_, err = db.Read(kg, id)
	assert.Error(t, err)

}
