package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Server struct {
	Name    string `json:"name"`
	Addr    string `json:"addr"`
	Version string `json:"version"`
	Weight  int    `json:"weight"`
	Ttl     int64  `json:"ttl"`
}

func (s Server) BuildRegisterKey() string {
	if s.Version == "" {
		return fmt.Sprintf("/%s/%s", s.Name, s.Addr)
	}
	return fmt.Sprintf("/%s/%s/%s", s.Name, s.Version, s.Addr)
}

func ParseValue(v []byte) (Server, error) {
	var server Server
	if err := json.Unmarshal(v, &server); err != nil {
		return server, err
	}
	return server, nil
}

func ParseKey(key string) (Server, error) {
	strs := strings.Split(key, "/")
	strs = filter(strs, func(s string) bool {
		return s != ""
	})
	if len(strs) == 2 {
		return Server{
			Name: strs[0],
			Addr: strs[1],
		}, nil
	}
	if len(strs) == 3 {
		return Server{
			Name:    strs[0],
			Addr:    strs[2],
			Version: strs[1],
		}, nil
	}
	return Server{}, errors.New("invalid key")
}

// filter函数用于过滤切片中的元素
func filter(strs []string, f func(string) bool) []string {
	var result []string
	for _, str := range strs {
		if f(str) {
			result = append(result, str)
		}
	}
	return result
}
