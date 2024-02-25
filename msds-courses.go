package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"
)

type MSDSCourse struct {
	CID string `json:"courseI_D`
	CNAME string `json:"course_name"`
	CPREREQ string `json:"prerequisite"` }

// JSONFILE resides in the current directory
var CSVFILE = "./coursedata.csv"

type MSDSCourseCatalog []MSDSCourse

var data = MSDSCourseCatalog{}
var index map[string]int

func readCSVFile(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil {
		return err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	// CSV file read all at once
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}
	
	// Getting the data out of the CSV and storing it in a temp variable 
	for _, line := range lines {
		temp := MSDSCourse{
			CID:       line[0],
			CNAME:    line[1],
			CPREREQ:        line[2],
		}
		// Storing data taken from the CSV to global variable which is a dictionary typed MSDSCourseCatalog
		data = append(data, temp)
	}

	return nil
}

func saveCSVFile(filepath string) error {
	// creating filepath variable methods on the returned File can be used for I/O
	csvfile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer csvfile.Close()

	// returns a new writer that writes to csvfile, and a writer writes records using CSV encoding
	csvwriter := csv.NewWriter(csvfile)
	for _, row := range data {
		temp := []string{row.CID, row.CNAME, row.CPREREQ}
		// writes a single CSV record from temp along with any necessary quoting
		_ = csvwriter.Write(temp)
	}
	csvwriter.Flush()
	return nil
}

func createIndex() error {
	index = make(map[string]int)
	for i, k := range data {
		// using the CID as the index
		key := k.CID
		index[key] = i
	}
	return nil
}

// Initialized by the user – returns a pointer
// If it returns nil, there was an error
func initS(ID, N, P string) *MSDSCourse {
	// Both of them should have a value
	if ID == "" || N == "" {
		return nil
	}
	// Creating a pointer for all the variables
	return &MSDSCourse{CID: ID, CNAME: N, CPREREQ: P}
}

func insert(pS *MSDSCourse) error {
	// If it already exists, do not add it
	_, ok := index[(*pS).CID]
	if ok {
		return fmt.Errorf("%s already exists", pS.CID)
	}

	data = append(data, *pS)
	// Update the index
	_ = createIndex()

	// saving the new csv file after the insert
	err := saveCSVFile(CSVFILE)
	if err != nil {
		return err
	}
	return nil
}

func deleteMSDSCourse(key string) error {
	i, ok := index[key]
	// if you search for a key that does not exist
	if !ok {
		return fmt.Errorf("%s cannot be found!", key)
	}
	data = append(data[:i], data[i+1:]...)
	// Update the index after a delete because key does not exist any more
	delete(index, key)

	err := saveCSVFile(CSVFILE)
	if err != nil {
		return err
	}
	return nil
}

func search(key string) *MSDSCourse {
	i, ok := index[key]
	if !ok {
		return nil
	}
	return &data[i]
}

func list() string {
	var all string
	for _, k := range data {
		all = all + k.CID + " " + k.CNAME + " " + k.CPREREQ + "\n"
	}
	return all
}

func main() {
	// read in the csv file
	err := readCSVFile(CSVFILE)
	if err != nil {
		fmt.Println(err)
		return
	}

	// create an index with the function based on the CID
	err = createIndex()
	if err != nil {
		fmt.Println("Cannot create index.")
		return
	}

	// Create a server with the http package
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         PORT,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	mux.Handle("/list", http.HandlerFunc(listHandler))
	mux.Handle("/insert/", http.HandlerFunc(insertHandler))
	mux.Handle("/insert", http.HandlerFunc(insertHandler))
	mux.Handle("/search", http.HandlerFunc(searchHandler))
	mux.Handle("/search/", http.HandlerFunc(searchHandler))
	mux.Handle("/delete/", http.HandlerFunc(deleteHandler))
	mux.Handle("/status", http.HandlerFunc(statusHandler))
	mux.Handle("/", http.HandlerFunc(defaultHandler))

	fmt.Println("Ready to serve at", PORT)
	err = s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
