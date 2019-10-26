package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v7"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, prs := pathsToUrls[r.URL.String()]
		if !prs {
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, url, http.StatusFound)
		}
	}
}

type pathMapping struct {
	Path string
	URL  string
}

func parseYAML(yml []byte) ([]pathMapping, error) {
	var mappings []pathMapping

	err := yaml.Unmarshal(yml, &mappings)
	if err != nil {
		return nil, err
	}

	return mappings, nil
}

func parseJSON(jsn []byte) ([]pathMapping, error) {
	var mappings []pathMapping

	err := json.Unmarshal(jsn, &mappings)
	if err != nil {
		return nil, err
	}

	return mappings, nil
}

func buildMap(mappings []pathMapping) map[string]string {
	pathMap := make(map[string]string)
	for _, m := range mappings {
		pathMap[m.Path] = m.URL
	}

	return pathMap
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	mappings, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(mappings)
	return MapHandler(pathMap, fallback), nil
}

// JSONHandler parses the provided JSON and returns an http.HandlerFunc
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	mappings, err := parseJSON(jsn)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(mappings)
	return MapHandler(pathMap, fallback), nil
}

// RedisHandler uses the provided Redis database to attempt
// to map any paths to their corresponding URL. If the path
// is not found in the database, then the fallback http.Handler
// will be called instead.
func RedisHandler(host string, port int, db int, key string, fallback http.Handler) http.HandlerFunc {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})

	return func(w http.ResponseWriter, r *http.Request) {
		url, err := client.HGet(key, r.URL.String()).Result()
		if err != nil {
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, url, http.StatusFound)
		}
	}
}
