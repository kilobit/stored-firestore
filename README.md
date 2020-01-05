stored-firestore
================

A GCP Firestore driver for the [StorEd](https://kilobit.ca/go/stored)
data respository.

Status: In-Development

Instantiate FireStore, then use the stored.Store interface.

```{.go}
// Using default marshal / unmarshal
var fs stored.Store = NewFireStore(project, nil, nil)

func Domain(store *stored.Store) {

	err := store.StoreItem(id, data)
	if err != nil {
		// Handle storage error.
	}

	obj, err := fs.Retrieve(id)
	if err != nil {
		// Handle read error.
	}
}

Domain(fs)
```

Features
--------

- Store data in the GCP cloud using Firestore.
- Decoupled domain and storage logic.

Upcoming:
- Transactions.
- Triggers.
- Limit operations by Criterria.

Installation
------------

```{.bash}
go get kilobit.cs/go/stored-firestore
```

Building
--------

```{.bash}
cd kilobit.ca/go/stored-firestore
go test -v
go build
```

Contribute
----------

Please submit a pull request with any bug fixes or feature requests
that you have.  All submissions imply consent to use / distribute
under the terms of the LICENSE.

Support
-------

Submit tickets through [github](https://github.com/kilobit/stored-firestore).

License
-------

See LICENSE.

Links to 3rd party content that may have it's own licensing.

--  
Created: Dec 17, 2019  
By: Christian Saunders <cps@kilobit.ca>  
Copyright 2019 Kilobit Labs Inc.  
