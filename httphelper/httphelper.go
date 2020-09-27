package httphelper

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Helper обертка для запроса/ответа
type Helper struct {
	log *logrus.Entry
	w   http.ResponseWriter
	r   *http.Request
}

// New создает обертку
func New(log *logrus.Entry, w http.ResponseWriter, r *http.Request) *Helper {
	return &Helper{
		log: log,
		w:   w,
		r:   r,
	}
}

// QueryValue возвращает параметр запроса
func (rw *Helper) QueryValue(name string) string {
	return rw.r.URL.Query().Get(name)
}

// Response отсылает ответ
func (rw *Helper) Response(code int, data []byte, headers Headers) {
	if headers != nil {
		for name, value := range headers {
			rw.w.Header().Add(name, value)
		}
	}

	rw.w.WriteHeader(code)

	if data != nil {
		rw.w.Write(data)
	}
}

// ResponseError отсылает ответ
func (rw *Helper) ResponseError(code int, err error) {
	rw.Response(code, []byte(err.Error()), nil)
}

// ResponseString отсылает ответ
func (rw *Helper) ResponseString(code int, data string) {
	rw.Response(code, []byte(data), nil)
}

// ResponseJSON отсылает ответ
func (rw *Helper) ResponseJSON(code int, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		rw.log.WithError(err).Error("Object parsing")
		rw.ResponseError(http.StatusInternalServerError, err)
		return
	}

	rw.Response(code, data, rw.Headers().Add("Content-Type", "application/json"))
}

// Headers заголовки
type Headers map[string]string

// Headers создает заголовки
func (rw *Helper) Headers() Headers {
	return make(Headers)
}

// Add добавляет заголовок
func (h Headers) Add(key, value string) Headers {
	h[key] = value
	return h
}
