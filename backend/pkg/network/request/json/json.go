package json

import (
	"blockchain/pkg/network"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func URLParamToInt(r *http.Request, key string) int {
	i, err := strconv.Atoi(chi.URLParam(r, key))
	if err != nil {
		return 0
	}
	return i
}

func URLParamToUint64(r *http.Request, key string) uint64 {
	value := chi.URLParam(r, key)
	ret, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0
	}
	return ret
}

func DefaultQuery(r *http.Request, key, def string) string {
	str := Query(r, key)
	if str == "" {
		return def
	}
	return str
}

func DefaultQueryBool(r *http.Request, key string, def bool) bool {
	str := Query(r, key)
	if str == "" {
		return def
	}
	return str == "true"
}

func DefaultQueryMapStringList(r *http.Request, key string, def map[string][]string) map[string][]string {
	str := Query(r, key)
	if str == "" {
		return def
	}

	var mapData map[string][]string
	if err := json.Unmarshal([]byte(str), &mapData); err != nil {
		return def
	}
	return mapData
}

func Query(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func DefaultQueryInt(r *http.Request, key string, def int) int {
	i, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil {
		return def
	}
	return i
}

func QueryUInt64(r *http.Request, key string) uint64 {
	value := r.URL.Query().Get(key)
	ret, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0
	}
	return ret
}

func Bind(out interface{}, r *http.Request) error {
	encoding := r.Header.Get("content-type")
	if strings.HasPrefix(encoding, network.ApplicationJson) {
		return json.NewDecoder(r.Body).Decode(out)
	}
	return errors.New("response: unexpected content type")
}

func GetFileContentType(r *http.Request, key string) string {
	file, header, err := r.FormFile(key)
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	if err != nil {
		return ""
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return ""
	}
	r.Body = io.NopCloser(bytes.NewBuffer(content))
	headerContentType := header.Header.Get("Content-Type")
	binaryContentType := http.DetectContentType(content)
	if headerContentType == binaryContentType {
		return binaryContentType
	}

	return headerContentType
}

func GetClientIP(r *http.Request) string {
	header := r.Header.Get("X-Forwarded-For")
	if header == "" {
		if strings.Contains(r.RemoteAddr, ":") {
			return strings.Split(r.RemoteAddr, ":")[0]
		}
		return r.RemoteAddr
	}
	return strings.Split(header, ",")[0]
}
