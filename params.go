package operator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

type paramsHolder struct {
	url   map[string]string
	query map[string][]string
	form  map[string][]string
	files map[string][]*multipart.FileHeader
}

func ParamsFrom(req *http.Request) map[string]string {
	// Parse the body depending on the content type.
	var (
		vals     map[string][]string
		files    map[string][]*multipart.FileHeader
		jsonVals map[string]string

		toReturn = make(map[string]string)
	)

	if contentTypeSlice, ok := req.Header["Content-Type"]; !ok || len(contentTypeSlice) < 1 {
		return toReturn
	}

	switch req.Header["Content-Type"][0] {
	case "application/x-www-form-urlencoded":
		// Typical form.
		if err := req.ParseForm(); err != nil {
			// TODO(ttacon): do something decent here
			return nil
		} else {
			vals = req.Form
		}

	case "multipart/form-data":
		// Multipart form.
		// TODO: Extract the multipart form param so app can set it.
		//  32 MB
		if err := req.ParseMultipartForm(32 << 20); err != nil {
			// TODO(ttacon): do something decent here
		} else {
			vals = req.MultipartForm.Value
			files = req.MultipartForm.File
		}
	case "application/json":
		// TODO(ttacon): check c.parseBody
		dataBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println("err : ", err)
			break
		}

		err = json.Unmarshal(dataBytes, &jsonVals)
		if err != nil {
			log.Println("err: ", err)
			break
		}

		for k, v := range jsonVals {
			toReturn[k] = v
		}
	}

	for k, v := range vals {
		toReturn[k] = v[0]
	}
	fmt.Println("val: ", vals)
	fmt.Println("files: ", files)

	return toReturn
}
