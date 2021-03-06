/* Copyright 2019 Kilobit Labs Inc. */

package firestore // import "kilobit.ca/go/stored-firestore"

import ctx "context"
import . "kilobit.ca/go/stored"
import "cloud.google.com/go/firestore"

//import "google.golang.org/api/iterator"
import "google.golang.org/api/option"

type Option func(*FireStore)

type Marshaler func(Storable) (interface{}, error)

func nopMarshaler(obj Storable) (interface{}, error) {
	return obj, nil
}

type UnMarshaler func(interface{}) (Storable, error)

func nopUnMarshaler(obj interface{}) (Storable, error) {
	return obj, nil
}

type FireStore struct {
	project     string                // Name of the GCP project
	client      *firestore.Client     // Client connection
	client_opts []option.ClientOption // Client options
	collection  string                // FS collection
	m           Marshaler             // Prepare Storables
	u           UnMarshaler           // Reconstitute Storables
}

func OptCollection(collection string) Option {
	return func(fs *FireStore) {
		fs.collection = collection
	}
}

func OptMarshaler(m Marshaler) Option {
	return func(fs *FireStore) {
		fs.m = m
	}
}

func OptUnMarshaler(u UnMarshaler) Option {
	return func(fs *FireStore) {
		fs.u = u
	}
}

func NewFireStore(project string, opts ...Option) *FireStore {

	fs := &FireStore{
		project,
		nil,
		[]option.ClientOption{},
		"",
		nopMarshaler,
		nopUnMarshaler,
	}

	fs.Options(opts...)

	return fs
}

func (fs *FireStore) Options(opts ...Option) {
	for _, opt := range opts {
		opt(fs)
	}
}

func (fs *FireStore) connect() error {
	if fs.client != nil {
		return nil
	}

	c, err := firestore.NewClient(
		ctx.TODO(),
		fs.project,
		fs.client_opts...)

	fs.client = c

	return err
}

func (fs *FireStore) setCollection(id ID) ID {
	if fs.collection != "" {
		id = (ID)(fs.collection + "/" + string(id))
	}

	return id
}

func (fs *FireStore) Close() {
	fs.client.Close()
	fs.client = nil
}

func (fs *FireStore) StoreItem(id ID, obj Storable) error {

	err := fs.connect()
	if err != nil {
		return err
	}

	id = fs.setCollection(id)

	mobj, err := fs.m(obj)
	if err != nil {
		return err
	}

	dr := fs.client.Doc((string)(id))
	_, err = dr.Set(ctx.TODO(), mobj)

	return err
}

func (fs *FireStore) Retrieve(id ID) (Storable, error) {

	err := fs.connect()
	if err != nil {
		return nil, err
	}

	id = fs.setCollection(id)

	dr := fs.client.Doc((string)(id))
	ds, err := dr.Get(ctx.TODO())
	if err != nil {
		return nil, err
	}

	obj, err := fs.u(ds.Data())
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// If the FS collection is set, return a pointer to that CollectionRef
// otherwise return pointers to all of the CollectionRefs in the FS.
//
func (fs *FireStore) listCollections() ([]*firestore.CollectionRef, error) {

	err := fs.connect()
	if err != nil {
		return nil, err
	}

	if fs.collection != "" {
		col := fs.client.Collection(fs.collection)
		return []*firestore.CollectionRef{col}, nil
	}

	return fs.client.Collections(ctx.TODO()).GetAll()
}

// Currently lists ids for all documents in the entire store.
//
func (fs *FireStore) List() ([]ID, error) {

	ids := []ID{}

	err := fs.connect()
	if err != nil {
		return nil, err
	}

	cols, err := fs.listCollections()
	if err != nil {
		return nil, err
	}

	for _, col := range cols {
		docs, err := col.DocumentRefs(ctx.TODO()).GetAll()
		if err != nil {
			return nil, err
		}

		for _, doc := range docs {

			id := doc.ID
			if fs.collection == "" {
				id = col.ID + "/" + id
			}

			ids = append(ids, (ID)(id))
		}
	}

	return ids, nil
}

func (fs *FireStore) Apply(f ItemHandler) error {

	ids, err := fs.List()
	if err != nil {
		return err
	}

	for _, id := range ids {
		obj, err := fs.Retrieve(id)
		if err != nil {
			return err
		}

		err = f(id, obj)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *FireStore) Delete(id ID) error {

	err := fs.connect()
	if err != nil {
		return err
	}

	id = fs.setCollection(id)

	dr := fs.client.Doc((string)(id))
	_, err = dr.Delete(ctx.TODO())

	return err
}
