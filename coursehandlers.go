package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const PORT = ":8080"


/*The logic behind each handler function is simple: read user input and 
decide whether the given input and HTTP method are the desired ones. 
The principles of each client interaction are also simple: the server 
should send appropriate error messages and HTTP codes back to the client 
so that everyone knows what really happened. Lastly, everything should be 
documented to communicate in a common language. */


// the default handler might print instructions about the operation of 
// the server as well as the list of available endpoints
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := "Thanks for visiting!\n"
	fmt.Fprintf(w, "%s", Body)
}

// The DELETE HTTP method looks like a rational choice when deleting a resource
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Get telephone
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)
	if len(paramStr) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found: "+r.URL.Path)
		return
	}

	log.Println("Serving:", r.URL.Path, "from", r.Host)


	CID := paramStr[2]
	// calling the delete function from other file to delete an entry 
	err := deleteMSDSCourse(CID)
	if err != nil {
		fmt.Println(err)
		Body := err.Error() + "\n"
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", Body)
		return
	}

	Body := CID + " deleted!\n"
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", Body)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := list()
	fmt.Fprintf(w, "%s", Body)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := fmt.Sprintf("Total entries: %d\n", len(data))
	fmt.Fprintf(w, "%s", Body)
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	// Split URL
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)

	if len(paramStr) < 5 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not enough arguments: "+r.URL.Path)
		return
	}

	courseI_D := paramStr[2]
	course_name := paramStr[3]
	prerequisite := paramStr[4]

	temp := &MSDSCourse{CID: courseI_D, CNAME: course_name, CPREREQ: prerequisite}
	err := insert(temp)

	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		Body := "Failed to add record\n"
		fmt.Fprintf(w, "%s", Body)
	} else {
		log.Println("Serving:", r.URL.Path, "from", r.Host)
		Body := "New record added successfully\n"
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", Body)
	}

	log.Println("Serving:", r.URL.Path, "from", r.Host)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Get Search value from URL
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)

	if len(paramStr) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found: "+r.URL.Path)
		return
	}

	var Body string
	CID := paramStr[2]
	t := search(CID)
	if t == nil {
		w.WriteHeader(http.StatusNotFound)
		Body = "Could not be found: " + CID + "\n"
	} else {
		w.WriteHeader(http.StatusOK)
		Body = t.CID + " " + t.CNAME + " " + t.CPREREQ + "\n"
	}

	fmt.Println("Serving:", r.URL.Path, "from", r.Host)
	fmt.Fprintf(w, "%s", Body)
}
