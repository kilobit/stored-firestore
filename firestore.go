/* Copyright 2019 Kilobit Labs Inc. */

package firestore // import "kilobit.ca/go/stored-firestore"

import ctx "context"
import . "kilobit.ca/go/stored"
import gfs "cloud.google.com/go/firestore"
//import "google.golang.org/api/iterator"
import "google.golang.org/api/option"

type Option func(*FireStore)

type Marshaler interface {
	Marshal(obj Storable) (interface{}, error)
}

type UnMarshaler interface {
	UnMarshal(obj interface{}) (Storable, error)
}

type FireStore struct {
	project string  // Name of the GCP project
	client *gfs.Client // Client connection
	client_opts []option.ClientOption // Client options
	m Marshaler // Prepare Storables
	u UnMarshaler // Reconstitute Storables
}

func NewFireStore(project string, m Marshaler,
	u UnMarshaler, opts... Option) *FireStore {

	fs := &FireStore{project, nil, []option.ClientOption{}, m, u}
	
	fs.Options(opts...)

	return fs
}

func (fs *FireStore) Options(opts... Option) {
	for _, opt := range opts {
		opt(fs)
	}
}

func (fs *FireStore) connect() error {
	if(fs.client != nil) {
		return nil
	}

	c, err := gfs.NewClient(
		ctx.TODO(),
		fs.project,
		fs.client_opts...)

	fs.client = c
	
	return err
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

	dr := fs.client.Doc((string)(id))
	_, err = dr.Set(ctx.TODO(), obj)
	
	return err
}

func (fs *FireStore) Retrieve(id ID) (Storable, error) {

	err := fs.connect()
	if err != nil {
		return nil, err
	}

	dr := fs.client.Doc((string)(id))
	ds, err := dr.Get(ctx.TODO())
	if err != nil {
		return nil, err
	}
	
	return ds.Data(), nil
}

func (fs *FireStore) List() ([]ID, error) {

	return []ID{}, nil
}

func (fs *FireStore) Apply(f ItemHandler) error {

	return nil
}

func (fs *FireStore) Delete(id ID) error {

	err := fs.connect()
	if err != nil {
		return err
	}

	dr := fs.client.Doc((string)(id))
	_, err = dr.Delete(ctx.TODO())
	
	return err 
}
