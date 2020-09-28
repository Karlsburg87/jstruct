package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//CLI options
var (
	data     = flag.String("data", "", "The raw json data. if \"-file\" or \"f\" flag is given then file contents will overwrite this")
	jsonFile string
	outfile  = flag.String("out", "newStruct.go", "Location to create the go file containing struct")
	name     = flag.String("name", "newStruct", "The name of the struct in the output file")
)

func main() {
	flag.StringVar(&jsonFile, "file", "", "The path to the file containing json data")
	flag.StringVar(&jsonFile, "f", "", "The path to the file containing json data")
	flag.Parse()

	//Ask for Json if nothing given
	if jsonFile == "" && *data == "" {
		log.Fatalln("Please feed data into the program by using the -file or -data flags")
	}

	//get data
	if jsonFile != "" {
		readFileBytes, err := ioutil.ReadFile(jsonFile)
		if err != nil && len(*data) <= 0 {
			log.Panicln("No data provided. Please ensure jsonFile path is correct.")
		}
		*data = string(readFileBytes)
	}

	//Unmarshal data into intermediatory map
	transition := make(map[string]interface{})

	err := json.Unmarshal([]byte(*data), &transition)
	if err != nil {
		log.Panicln(err)
	}

	//Print the struct
	f, err := os.Create(*outfile)
	defer f.Close()

	//Opening document information
	fmt.Fprintf(f, "package %s", *name)

	//Create the new main struct
	splitOuts := PrintStruct(f, transition, *name)

	//Now create the new structs where necarssary
	worker := splitOuts

	fieldsWithStructValuesToCreate := make([]interface{}, 0) //!error check
	for {
		if len(worker) <= 0 {
			break
		}

		worker2 := make(map[string]interface{})
		for ke := range worker {
			fval, ok := worker[ke].(map[string]interface{})
			if ok {
				worker2 = PrintStruct(f, fval, ke)
			} else {
				fieldsWithStructValuesToCreate = append(fieldsWithStructValuesToCreate, ke)
			}
		}

		//Make sure we only create a sub struct once
		for k, v := range worker2 {
			if splitOuts[k] != nil {
				delete(worker, k)
			} else {
				splitOuts[k] = v
			}
		}
		worker = worker2
	}

	//Print issues
	if len(fieldsWithStructValuesToCreate) > 0 {
		log.Printf("%d issues: %v\n", len(fieldsWithStructValuesToCreate), fieldsWithStructValuesToCreate)
	}
	/*********************
	Sign-off
	*********************/
	fmt.Fprintln(f, "\n\n/***********\nA Rhythmic Sound Project\n***********/")

}

//PrintStruct takes a map and outputs the struct to file. Returns map of sub structs to create
func PrintStruct(f *os.File, m map[string]interface{}, structName string) map[string]interface{} {
	//For sub structs
	collection := make(map[string]interface{})

	//Proper naming convention imposing
	structName = NamingConvention(structName)

	//Open struct
	fmt.Fprintf(f, " \n\ntype %s struct {\n", structName)

	//Print main body to file
	for key, value := range m {
		newKey := NamingConvention(key)

		switch v := value.(type) {
		case map[string]interface{}:
			//add to collection to form own structs in output
			collection[structName+newKey] = v
			//add field with type of to be created struct here
			fmt.Fprintf(f, "\t%s %s `json:\"%s,omitempty\"`\n", newKey, structName+newKey, key)

		case []map[string]interface{}:
			//add to collection to form own structs in output
			collection[structName+newKey] = v[0]
			//add field with type of to be created struct here
			fmt.Fprintf(f, "\t%s []%s `json:\"%s,omitempty\"`\n", newKey, structName+newKey, key)

		case []interface{}:
			switch sv := v[0].(type) {
			case map[string]interface{}:
				//add to collection to form own structs in output
				collection[structName+newKey] = sv
				//add field with type of to be created struct here
				fmt.Fprintf(f, "\t%s []%s `json:\"%s,omitempty\"`\n", newKey, structName+newKey, key)

			default:
				fmt.Fprintf(f, "\t%s []%T `json:\"%s,omitempty\"`\n", newKey, sv, key)
			}

		default:
			fmt.Fprintf(f, "\t%s %T `json:\"%s,omitempty\"`\n", newKey, value, key)
		}
	}

	//Close struct bracket
	fmt.Fprintln(f, "}")

	return collection
}

//NamingConvention takes for structs and imposes a singular naming convention
func NamingConvention(nameRaw string) string {
	c := regexp.MustCompile(`[-_](\w?)`)
	newKey := c.ReplaceAllStringFunc(nameRaw, func(s string) string { return strings.ToUpper(s[len(s)-1:]) }) //CamelCase
	newKey = strings.ToUpper(newKey[:1]) + newKey[1:]                                                         //Capitalise
	if strings.Contains(newKey[len(newKey)-2:], "id") {                                                       //id to ID
		newKey = newKey[:len(newKey)-2] + "ID"
	}
	newKey = strings.ReplaceAll(newKey, "url", "URL")          //url to URL
	if _, err := strconv.Atoi(string(newKey[0])); err == nil { //Does not start with number
		newKey = "N" + newKey
	}
	//newKey = regexp.MustCompile(`[\[\]\-_\\/!"£"$%^&()*+=\:;'#<>?,.|]`).ReplaceAllString(newKey, "") //remove remaining non alphanumeric charecters
	return newKey
}
