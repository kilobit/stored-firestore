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

	fs := NewFireStore(project, OptCollection("testing-axyz7b"))

	return fs
}

var SNRData map[ID]Storable = map[ID]Storable{
	"test1": map[string]interface{}{
		"foo":  "bar",
		"bing": "bong",
	},

	"test2": map[string]interface{}{
		"foo": 12,
		"bing": map[string]interface{}{
			"nested": true,
		},
	},

	"test3": map[string]interface{}{
		"foo":  42,
		"bing": 43,
	},
}

// This test should probably be broken up but currently serves the
// need of a quick implementation that demonstrates the golden path
// functionality.
//
func TestFireStoreStoreAndRetrieve(t *testing.T) {

	fs := newTestFireStore(t)
	defer fs.Close()

	// Store the test objects.
	for id, data := range SNRData {

		err := fs.StoreItem(id, data)
		if err != nil {
			t.Errorf("%v\n%s", id, err)
		}
	}

	// List the test objects
	ids, err := fs.List()
	if err != nil {
		t.Error(err)
	}

	for _, id := range ids {
		_, ok := SNRData[id]
		if !ok {
			t.Errorf("Missing test data %s from ids.\n", id)
		}
	}

	// Test Apply.
	err = fs.Apply(func(id ID, obj Storable) error {
		data, ok := SNRData[id]
		if !ok {
			t.Errorf("ID %s in Apply not associated to test data.", id)

		}

		assert.ExpectDeep(t, data, obj)

		return nil
	})
	if err != nil {
		t.Error(err)
	}

	// Retrieve and compare the test objects.
	for id, data := range SNRData {
		obj, err := fs.Retrieve(id)
		if err != nil {
			t.Errorf("%v\n%s", id, err)
		}

		// Useful for debugging.
		// t.Logf("%#v", data)
		// t.Logf("%#v", obj)

		assert.ExpectDeep(t, data, obj)
	}

	// Delete the test objects.
	for id, _ := range SNRData {

		err := fs.Delete(id)
		if err != nil {
			t.Error(err)
		}
	}
}
