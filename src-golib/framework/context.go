package framework

import (
	stdContext "context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lazygo/client/utils"
)

type Context interface {
	stdContext.Context

	Request() *Request
	// ResponseWriter 返回ResponseWriter
	ResponseWriter() *ResponseWriter

	QueryParam(name string) string
	Path() string
	GetBody() (io.ReadCloser, error)
	Bind(BaseRequest) error

	// WithValue 存入数据到当前请求的context
	WithValue(key string, val interface{})

	Stream(contentType string, r io.Reader) error
	Blob(contentType string, b []byte) error
	JSON(i interface{}) error
	Error(err error)

	Entry() int
	// IsDebug return the Server is debug.
	IsDebug() bool

	Logger() Logger
	// RequestId 获取请求id
	RequestID() uint64
	Succ(data interface{}) error
}

const (
	ENTRY_HTTP = iota
	ENTRY_DYLIB
)

type context struct {
	responseWriter *ResponseWriter
	lock           sync.RWMutex
	store          Map
	app            *Application
	request        *Request
}

func (c *context) Request() *Request {
	return c.request
}

func (c *context) ResponseWriter() *ResponseWriter {
	return c.responseWriter
}

func (c *context) WithValue(key string, val interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.store == nil {
		c.store = make(Map)
	}
	if val == nil {
		delete(c.store, key)
		return
	}
	c.store[key] = val
}

func (c *context) JSON(i interface{}) error {
	c.writeContentType(MIMEApplicationJSONCharsetUTF8)
	enc := json.NewEncoder(c.responseWriter)
	return enc.Encode(i)
}

func (c *context) Error(err error) {
	c.app.ErrorHandler(err, c)
}

func (c *context) IsDebug() bool {
	return c.app.Debug
}

func (c *context) QueryParam(name string) string {
	if c.request.URL == nil {
		return ""
	}
	return c.request.URL.Query().Get(name)
}

func (c *context) Path() string {
	if c.request.URL == nil {
		return ""
	}
	return c.request.URL.Path
}

func (c *context) GetBody() (io.ReadCloser, error) {
	return io.NopCloser(c.request.Data), nil
}

// Deadline returns that there is no deadline (ok==false) when c.Request has no Context.
func (c *context) Deadline() (deadline time.Time, ok bool) {
	return c.request.Context().Deadline()
}

// Done returns nil (chan which will wait forever) when c.Request has no Context.
func (c *context) Done() <-chan struct{} {
	return c.request.Context().Done()
}

// Err returns nil when c.Request has no Context.
func (c *context) Err() error {
	return c.request.Context().Err()
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *context) Value(key interface{}) interface{} {
	if keyAsString, ok := key.(string); ok {
		c.lock.RLock()
		val, ok := c.store[keyAsString]
		c.lock.RUnlock()
		if ok {
			return val
		}
	}
	return c.request.Context().Value(key)
}

func (c *context) Bind(v BaseRequest) error {
	// result pointer value
	rpv := reflect.ValueOf(v)
	if rpv.Kind() != reflect.Ptr || rpv.IsNil() {
		c.Logger().Error("[error][msg: bind value not a pointer]")
		return ErrInternalServerError
	}

	body, err := c.GetBody()
	if err != nil {
		return err
	}
	err = json.NewDecoder(body).Decode(v)
	if err != nil {
		return err
	}

	// result value
	rv := rpv.Elem()
	if rv.Kind() != reflect.Struct {
		c.Logger().Error("[msg: bind value not a struct pointer]")
		return ErrInternalServerError
	}

	for i := 0; i < rv.NumField(); i++ {
		if !rv.Field(i).CanSet() {
			continue
		}
		tField := rv.Type().Field(i)
		field := tField.Tag.Get("json")
		if field == "" {
			continue
		}

		binds := strings.Split(tField.Tag.Get("bind"), ",")
		var val interface{}
		for _, bind := range binds {
			switch bind {
			case "value":
				val = c.Value(field)
			case "query":
				val = c.QueryParam(field)
			default:
				continue
			}
			if val != "" && val != nil {
				break
			}
		}
		if val == nil || val == "" {
			continue
		}

		procList := strings.Split(tField.Tag.Get("process"), ",")
		if to, ok := toType(val, tField.Type, procList); ok {
			rv.Field(i).Set(reflect.ValueOf(to))
		}

	}
	return nil
}

// Logger 日志记录器
func (c *context) Logger() Logger {
	return &loggerImpl{
		Ctx: c,
	}
}

func (c *context) writeContentType(value string) {
	if httpResponseWriter, ok := c.responseWriter.Writer.(http.ResponseWriter); ok {
		header := httpResponseWriter.Header()
		if header.Get(HeaderContentType) == "" {
			header.Set(HeaderContentType, value)
		}
	}
}

func (c *context) Blob(contentType string, b []byte) error {
	c.writeContentType(contentType)
	_, err := c.responseWriter.Write(b)
	return err
}

func (c *context) Stream(contentType string, r io.Reader) error {
	c.writeContentType(contentType)
	_, err := io.Copy(c.responseWriter, r)
	return err
}

func (c *context) Entry() int {
	_, ok := c.responseWriter.Writer.(http.ResponseWriter)
	if ok {
		return ENTRY_HTTP
	}
	return ENTRY_DYLIB
}

// GetRequestID 获取请求id
func (c *context) RequestID() uint64 {
	rid, _ := c.Value(HeaderXRequestID).(uint64)
	return rid
}

// Succ 返回成功
func (c *context) Succ(data interface{}) error {
	result := Map{
		"code":  200,
		"errno": 0,
		"msg":   "ok",
		"data":  data,
		"rid":   c.RequestID(),
		"t":     time.Now().Unix(),
	}
	return c.JSON(result)
}

func toType(val interface{}, rType reflect.Type, procList []string) (interface{}, bool) {
	typeName := rType.String()
	typeKind := rType.Kind()

	if typeKind == reflect.Ptr {
		typeKind = rType.Elem().Kind()
	}
	rv := reflect.ValueOf(val)

	if (typeKind == reflect.Array || typeKind == reflect.Interface || typeKind == reflect.Map || typeKind == reflect.Slice || typeKind == reflect.Struct) && rv.Kind() == reflect.String {
		returnVal := reflect.New(rType)
		strVal := utils.ToString(val)
		if err := json.Unmarshal([]byte(strVal), returnVal.Interface()); err == nil {
			val = returnVal.Elem().Interface()
			rv = returnVal.Elem()
		}
	}

	if typeName == "interface {}" {
		return val, true
	}
	switch typeName {
	case "int":
		return utils.ToInt(val), true
	case "[]int":
		returnVal := make([]int, 0)
		strVal := utils.ToString(val)
		if len(strVal) == 0 {
			return returnVal, true
		}
		for _, str := range strings.Split(strVal, ",") {
			returnVal = append(returnVal, utils.ToInt(str))
		}
		return returnVal, true
	case "int8":
		return int8(utils.ToInt(val)), true
	case "[]int8":
		returnVal := make([]int8, 0)
		strVal := utils.ToString(val)
		if len(strVal) == 0 {
			return returnVal, true
		}
		for _, str := range strings.Split(strVal, ",") {
			returnVal = append(returnVal, int8(utils.ToInt(str)))
		}
		return returnVal, true
	case "int16":
		return int16(utils.ToInt(val)), true
	case "[]int16":
		returnVal := make([]int16, 0)
		strVal := utils.ToString(val)
		if len(strVal) == 0 {
			return returnVal, true
		}
		for _, str := range strings.Split(strVal, ",") {
			returnVal = append(returnVal, int16(utils.ToInt(str)))
		}
		return returnVal, true
	case "int32":
		return int32(utils.ToInt(val)), true
	case "[]int32":
		returnVal := make([]int32, 0)
		strVal := utils.ToString(val)
		if len(strVal) == 0 {
			return returnVal, true
		}
		for _, str := range strings.Split(strVal, ",") {
			returnVal = append(returnVal, int32(utils.ToInt(str)))
		}
		return returnVal, true
	case "int64":
		return utils.ToInt64(val), true
	case "[]int64":
		returnVal := make([]int64, 0)
		strVal := utils.ToString(val)
		if len(strVal) == 0 {
			return returnVal, true
		}
		for _, str := range strings.Split(strVal, ",") {
			returnVal = append(returnVal, utils.ToInt64(str))
		}
		return returnVal, true
	case "uint":
		return utils.ToUint(val), true
	case "[]uint":
		returnVal := make([]uint, 0)
		strVal := utils.ToString(val)
		if len(strVal) == 0 {
			return returnVal, true
		}
		for _, str := range strings.Split(strVal, ",") {
			returnVal = append(returnVal, utils.ToUint(str))
		}
		return returnVal, true
	case "uint8":
		return uint8(utils.ToUint(val)), true
	case "[]uint8":
		returnVal := make([]uint8, 0)
		strVal := utils.ToString(val)
		if len(strVal) == 0 {
			return returnVal, true
		}
		for _, str := range strings.Split(strVal, ",") {
			returnVal = append(returnVal, uint8(utils.ToUint(str)))
		}
		return returnVal, true
	case "uint16":
		return uint16(utils.ToUint(val)), true
	case "[]uint16":
		returnVal := make([]uint16, 0)
		strVal := utils.ToString(val)
		if len(strVal) == 0 {
			return returnVal, true
		}
		for _, str := range strings.Split(strVal, ",") {
			returnVal = append(returnVal, uint16(utils.ToUint(str)))
		}
		return returnVal, true
	case "uint32":
		return uint32(utils.ToUint(val)), true
	case "[]uint32":
		returnVal := make([]uint32, 0)
		strVal := utils.ToString(val)
		if len(strVal) == 0 {
			return returnVal, true
		}
		for _, str := range strings.Split(strVal, ",") {
			returnVal = append(returnVal, uint32(utils.ToUint(str)))
		}
		return returnVal, true
	case "uint64":
		return utils.ToUint64(val), true
	case "[]uint64":
		returnVal := make([]uint64, 0)
		strVal := utils.ToString(val)
		if len(strVal) == 0 {
			return returnVal, true
		}
		for _, str := range strings.Split(strVal, ",") {
			returnVal = append(returnVal, utils.ToUint64(str))
		}
		return returnVal, true
	case "float32":
		return float32(utils.ToFloat(val)), true
	case "float64":
		return utils.ToFloat(val), true
	case "string":
		return process(utils.ToString(val), procList), true
	case "[]string":
		returnVal := make([]string, 0)
		strVal := utils.ToString(val)
		if len(strVal) == 0 {
			return returnVal, true
		}
		for _, str := range strings.Split(strVal, ",") {
			returnVal = append(returnVal, process(utils.ToString(str), procList))
		}
	default:
		valType := rv.Type().String()
		if valType == typeName {
			return val, true
		}
		if strings.HasPrefix(valType, "*") && valType[1:] == typeName {
			val = rv.Elem().Interface()
			return val, true
		}
	}
	return val, false
}

func process(str string, procList []string) string {
	for _, proc := range procList {
		switch {
		case strings.HasPrefix(proc, "trim"):
			str = strings.TrimSpace(str)
		case strings.HasPrefix(proc, "tolower"):
			str = strings.ToLower(str)
		case strings.HasPrefix(proc, "toupper"):
			str = strings.ToUpper(str)
		case strings.HasPrefix(proc, "cut("):
			if n, err := strconv.Atoi(proc[4 : len(proc)-1]); err != nil {
				str = utils.CutRune(str, n)
			}
		}
	}
	return str
}
