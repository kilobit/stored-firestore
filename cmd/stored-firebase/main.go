/* Copyright 2019 Kilobit Labs Inc. */

package main

import "fmt"
import "io"
import "os"
import "kilobit.ca/go/args"
import "bufio"
import "encoding/json"

import . "kilobit.ca/go/stored"
import . "stored-firestore"

const PROJECT_ENV_NAME string = "GOOGLE_PROJECT_NAME"

const USAGE_MSG string = `
usage: stored-firestore [Options] Command ARGS...

Options:
-p, --project PROJECT_NAME
        Set or override the GCP project name in the environment.

Commands:

help    Print this usage message.

del ID  Delete the document referred to by ID.
get ID  Get the document referred to by ID.
set ID  Set the contents of stdin to the document referred to by
        ID.

Project Name:

Every GCP Firestore repository must have a related project name.
Either Set this name as the first argument preceeding the COMMAND or
in the environment varialbe called %s.

Authorization:

Be sure that authorization is provided by setting the
GOOGLE_APPLICATION_CREDENTIALS environment appropriately.

`

// Return a set of global options
//
func getGlobalOpts(ap *args.ArgParser) map[string]string {

	opts := map[string]string{}

	for {
		opt := ap.NextOpt()

		if opt == "" {
			break
		}

		switch opt {
		case "p", "project":
			opts["project"] = ap.NextArg()

		case "v", "verbose":
			opts["verbose"] = ""
		}
	}

	return opts
}

// Looks for the GCP project name in the options.  Failing that it
// will look in the environment and if still unsuccessful will error
// and quit the program.
//
func getProjectName(opts map[string]string) string {

	project, ok := opts["project"]
	if !ok {
		project, ok = os.LookupEnv(PROJECT_ENV_NAME)
		if !ok {
			fmt.Fprintln(os.Stderr, "Missing project name.")
			os.Exit(1)
		}
	}

	return project
}

// Write the usage message.
func help(w io.Writer) {

	_, err := fmt.Fprintf(w, USAGE_MSG, PROJECT_ENV_NAME)
	if err != nil {
		panic("Failed to write usage message.")
	}
}

// Get the document from the store with the specified id and write it
// to the given Writer.
//
// Returns an appropriate exit code.
//
func get(store Store, id ID, w io.Writer) int {

	obj, err := store.Retrieve(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	_, err = fmt.Fprintln(w, obj)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

// Set the document with the specified id and write it to the store.
//
// Returns an appropriate exit code.
//
func set(store Store, id ID, r io.Reader) int {

	br := bufio.NewReader(r)
	bs, err := br.ReadBytes('\000')
	if err != io.EOF {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	var obj interface{}
	err = json.Unmarshal(bs, &obj)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	err = store.StoreItem(id, &obj)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

// Delete the document with the given ID from the store.
//
// Returns an appropriate exit code.
//
func del(store Store, id ID) int {

	err := store.Delete(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func main() {

	ap := args.NewArgParser(os.Args[1:])

	opts := getGlobalOpts(ap)

	project := getProjectName(opts)

	cmd := ap.NextArg()

	switch cmd {

	case "del", "delete":
		id := (ID)(ap.NextArg())
		store := NewFireStore(project, nil, nil)
		code := del(store, id)
		os.Exit(code)

	case "set":
		id := (ID)(ap.NextArg())
		store := NewFireStore(project, nil, nil)
		code := set(store, id, os.Stdin)
		os.Exit(code)

	case "get":
		id := (ID)(ap.NextArg())
		store := NewFireStore(project, nil, nil)
		code := get(store, id, os.Stdout)
		os.Exit(code)

	case "help":
		help(os.Stdout)
		os.Exit(0)

	default:
		fmt.Fprintln(os.Stderr, "Unrecognized command, ", cmd)
		os.Exit(1)
	}
}
