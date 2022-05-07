package stringutil

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	numBytes    = "1234567890"
)

func ConvertBool(val string) bool {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return b
}

func ConvertInt(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}

func ConvertUint64ToString(val uint64) string {
	return strconv.FormatUint(val, 10)
}

func ConvertInt64ToString(val int64) string {
	return strconv.FormatInt(val, 10)
}

func ConvertInt32ToString(val int32) string {
	return strconv.FormatInt(int64(val), 10)
}

func ConvertStringToUint64(val string) uint64 {
	value, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return value
}

func ConvertStringToUint32(val string) uint32 {
	value, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(value)
}

func ConvertStringToInt64(val string) int64 {
	value, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0
	}
	return value
}

func ConvertDuration(val string) time.Duration {
	i, err := strconv.Atoi(val)
	if err != nil {
		return time.Duration(0)
	}
	return time.Duration(i)
}
func ConvertFloat64(val string) float64 {
	value, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0
	}
	return value
}

func RemoveLine(str string) string {
	return strings.NewReplacer(
		"\r\n", "",
		"\r", "",
		"\n", "",
	).Replace(str)
}

func ToJson(obj interface{}) string {
	var buf bytes.Buffer
	b, err := json.Marshal(obj)
	if err != nil {
		log.Printf("error: %v", err)
	}
	buf.Write(b)
	return buf.String()
}

func RandStr(n int) string {
	buf := make([]byte, n)
	max := new(big.Int)
	max.SetInt64(int64(len(letterBytes)))
	for i := range buf {
		r, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		buf[i] = letterBytes[r.Int64()]
	}
	return string(buf)
}

// RandStrNum digitで指定された桁数の数字文字列を返す
func RandStrNum(digit int) string {
	buf := make([]byte, digit)
	max := new(big.Int)
	max.SetInt64(int64(len(numBytes)))
	for i := range buf {
		r, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		buf[i] = numBytes[r.Int64()]
	}
	return string(buf)
}

// SliceCharByte Byte数でのSlice
func SliceCharByte(value string, limit int) string {
	if value == "" {
		return ""
	} else if len(value) <= limit {
		return value
	}
	return value[0:limit]
}

// SliceString 全角考慮した文字数
func SliceString(value string, limit int) string {
	if value == "" {
		return ""
	} else if len(value) <= limit {
		return value
	}
	return string([]rune(value)[:limit])
}

var htmlTagPattern = regexp.MustCompile(`(</?[a-zA-Z]+?[^>]*/?>)*`)

func RemoveHtmlTag(in string) string {
	groups := htmlTagPattern.FindAllString(in, -1)
	// should replace long string first
	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) > len(groups[j])
	})
	for _, group := range groups {
		if strings.TrimSpace(group) != "" {
			in = strings.ReplaceAll(in, group, "")
		}
	}
	return in
}

func StringOrDefault(ptr *string, def string) string {
	if ptr == nil {
		return def
	}
	return *ptr
}
