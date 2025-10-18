package main

import (
	"fmt"

	"github.com/aholstenson/loro-go"
)

func main() {
	doc := loro.NewLoroDoc()

	// Insert a value into a root map named "test"
	loroMap := doc.GetMap(loro.AsContainerId("test"))
	err := loroMap.Insert("key", loro.AsStringValue("value"))
	if err != nil {
		fmt.Println("Error inserting value:", err)
		return
	}

	// Export the document
	export, err := doc.ExportSnapshot()
	if err != nil {
		fmt.Println("Error exporting document:", err)
		return
	}

	fmt.Println("Exported document length:", len(export))
}
