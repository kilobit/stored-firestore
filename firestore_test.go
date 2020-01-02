/* Copyright 2019 Kilobit Labs Inc. */

package firestore // import "kilobit.ca/go/stored-firestore"

import . "kilobit.ca/go/stored"
import "kilobit.ca/go/tested/assert"
import "testing"
import "os"

const PROJECT_ENV_NAME string = "GOOGLE_PROJECT_NAME"

func TestFireStoreTest(t *testing.T) {
	assert.Expect(t, true, true)
}

func newTestFireStore(t *testing.T) *FireStore {
	
	project, ok := os.LookupEnv(PROJECT_ENV_NAME)
	if !ok {
		t.Skip(PROJECT_ENV_NAME + " environment variable not set.")
	}
	
	fs := NewFireStore(project, nil, nil)

	return fs
}

var SNRData map[ID]Storable = map[ID]Storable{
	"testing/test1": map[string]interface{}{
		"foo": "bar",
		"bing": "bong",
	},

	"testing/test2": map[string]interface{}{
		"foo": 12,
		"bing": map[string]interface{}{
			"nested": true,
		},
	},

	"testing/test3": map[string]interface{}{
		"foo": 42,
		"bing": 43,
	},
}

func TestFireStoreStoreAndRetrieve(t *testing.T) {
	
	fs := newTestFireStore(t)
	defer fs.Close()

	for id, data := range SNRData {
		
		err := fs.StoreItem(id, data)
		if err != nil {
			t.Errorf("%v\n%s", id, err)
		}

		obj, err := fs.Retrieve(id)
		if err != nil {
			t.Errorf("%v\n%s", id, err)
		}

		// Useful for debugging.
		// t.Logf("%#v", data)
		// t.Logf("%#v", obj)
		
		assert.ExpectDeep(t, data, obj)
	}

	for id, _ := range SNRData {

		err := fs.Delete(id)
		if err != nil {
			t.Error(err)
		}
	}
}
