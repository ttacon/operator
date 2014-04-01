package yolo

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Conductor struct {
	router *router
}

func NewConductor() *Conductor {
	return &Conductor{
		router: newRouter(),
	}
}

func (c *Conductor) Get(url string, handlers ...Handler) {
	c.router.addRoute("/GET"+url, handlers)
}

func (c *Conductor) Head(url string, handlers ...Handler) {
	c.router.addRoute("/HEAD"+url, handlers)
}

func (c *Conductor) Put(url string, handlers ...Handler) {
	c.router.addRoute("/PUT"+url, handlers)
}

func (c *Conductor) Post(url string, handlers ...Handler) {
	c.router.addRoute("/POST"+url, handlers)
}

func (c *Conductor) Delete(url string, handlers ...Handler) {
	c.router.addRoute("/DELETE"+url, handlers)
}

func (c *Conductor) Options(url string, handlers ...Handler) {
	c.router.addRoute("/OPTIONS"+url, handlers)
}

func (c *Conductor) Trace(url string, handlers ...Handler) {
	c.router.addRoute("/TRACE"+url, handlers)
}

func (c *Conductor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlers, urlParams := c.router.findRoute(r)
	if len(handlers) == 0 {
		http.NotFound(w, r)
		return
	}

	// get other params, form/body/etc...
	// 	params := make(map[string]interface{})

	// should we give an option to reuse pointered structs across functions?
	// what about returning values to be reused in next handler?

	for _, handler := range handlers {
		// get params for handler
		args, err := paramsFor(handler, w, r, urlParams)
		if err != nil {
			// TODO(ttacon): need better stuff to do than break...
			// 500 or 404?
			break
		}

		// call handler
		// why does the following not work? :
		//		ty := reflect.TypeOf(handler)
		//		h := reflect.New(ty)
		//		e := h.Elem()
		e := reflect.ValueOf(handler)
		e.Call(args)

		// can we check header map to see if written to?
		if len(w.Header()) > 0 {
			break
		}
	}
}

func paramsFor(h Handler, w http.ResponseWriter, r *http.Request, params map[string]interface{}) ([]reflect.Value, error) {
	ty := reflect.TypeOf(h)
	if ty.Kind() != reflect.Func {
		return nil, fmt.Errorf("%v is not of reflect.Kind Func", ty)
	}

	numArgs := ty.NumIn()
	vals := make([]reflect.Value, numArgs)
	for i := 0; i < numArgs; i++ {
		t := ty.In(i)
		if !validParam(reflect.New(t).Elem()) {
			return nil, fmt.Errorf("params to handler must be struct, "+
				"map, http.Request or http.ResponseWriter, type was: %v",
				t)
		}

		if t == reflect.TypeOf(HttpRequestType) {
			vals[i] = reflect.ValueOf(r)
			continue
		}

		if t.String() == "http.ResponseWriter" {
			vals[i] = reflect.ValueOf(w)
			continue
		}

		val := reflect.New(t)
		if t.Kind() == reflect.Map {
			// add everything that hasn't been used
			continue
		}

		e := val.Elem()
		for j := 0; j < t.NumField(); j++ {
			f := e.Field(j)
			param, ok := params[t.Field(j).Name]
			if !ok {
				param, ok = params[strings.ToLower(t.Field(j).Name)]
			}

			if !ok {
				// param doesn't exist in url, body, etc...
				continue
			}

			err := setField(f, param)
			if err != nil {
				return nil, err
			}
		}

		vals[i] = e
	}

	return vals, nil
}

func setField(field reflect.Value, param interface{}) error {
	// TODO(ttacon: allow for struct annotations to be used
	switch field.Kind() {
	case reflect.Int:
		fallthrough
	case reflect.Int64:
		// check if hex or oct, etc
		// how to know bit size? <--- I guess default to 64 since
		// otherwise type would be intDD
		i, err := strconv.ParseInt(param.(string), 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Int32:
		// TODO(ttacon): deal with hex/oct/binary/etc...
		i, err := strconv.ParseInt(param.(string), 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Int16:
		// TODO(ttacon): deal with hex/oct/binary/etc...
		i, err := strconv.ParseInt(param.(string), 10, 16)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Int8:
		// TODO(ttacon): deal with hex/oct/binary/etc...
		i, err := strconv.ParseInt(param.(string), 10, 8)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint:
		fallthrough
	case reflect.Uint64:
		// check if hex or oct, etc
		// how to know bit size? <--- I guess default to 64 since
		// otherwise type would be intDD
		i, err := strconv.ParseUint(param.(string), 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.Uint32:
		// TODO(ttacon): deal with hex/oct/binary/etc...
		i, err := strconv.ParseUint(param.(string), 10, 32)
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.Uint16:
		// TODO(ttacon): deal with hex/oct/binary/etc...
		i, err := strconv.ParseUint(param.(string), 10, 16)
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.Uint8:
		// TODO(ttacon): deal with hex/oct/binary/etc...
		i, err := strconv.ParseUint(param.(string), 10, 8)
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(param.(string))
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Float32:
		d, err := strconv.ParseFloat(param.(string), 32)
		if err != nil {
			return err
		}
		field.SetFloat(d)
	case reflect.Float64:
		d, err := strconv.ParseFloat(param.(string), 64)
		if err != nil {
			return err
		}
		field.SetFloat(d)
	case reflect.String:
		field.SetString(param.(string))
	default:
		return fmt.Errorf(
			"tried to set unsupported value: %v, of type: %v",
			param,
			field.Type())
	}

	return nil
}

var (
	HttpRequestType    *http.Request
	HttpResponseWriter http.ResponseWriter
)

func validParam(t reflect.Value) bool {
	if t.Kind() == reflect.Struct {
		return true
	}

	if t.Type().String() == "http.ResponseWriter" {
		return true
	}

	if t.Type() == reflect.TypeOf(HttpRequestType) {
		return true
	}

	// ensure is map[string]interface{}
	return false
}
