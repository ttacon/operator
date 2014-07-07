package operator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

type paramsHolder struct {
	url   map[string]string
	query map[string][]string
	form  map[string][]string
	files map[string][]*multipart.FileHeader
}

func ParamsFrom(req *http.Request) map[string]interface{} {
	// Parse the body depending on the content type.
	var (
		vals     map[string][]string
		files    map[string][]*multipart.FileHeader
		jsonVals map[string]interface{}

		toReturn = make(map[string]interface{})
	)

	if contentTypeSlice, ok := req.Header["Content-Type"]; !ok || len(contentTypeSlice) < 1 {
		return toReturn
	}

	fmt.Println(req.Header["Content-Type"][0])
	fmt.Println(req.Header["Content-Type"])
	contentType := req.Header["Content-Type"][0]
	if strings.HasPrefix(contentType, "multipart/form-data") {
		contentType = "multipart/form-data"
	}

	switch contentType {
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

		// crappers...
		// TODO(ttacon): jsonVals should be map[string]interface{}
		err = json.Unmarshal(dataBytes, &jsonVals)
		if err != nil {
			log.Println("err: ", err)
			break
		}

		for k, v := range jsonVals {
			toReturn[k] = v
		}
	}
	// TODO(ttacon): xml - lulz, people still use this?

	for k, v := range vals {
		toReturn[k] = v[0]
	}
	fmt.Println("val: ", vals)
	fmt.Println("files: ", files)

	return toReturn
}
