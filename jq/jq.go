package jq

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/itchyny/gojq"
)

const debug = false

//go:embed *
var scripts embed.FS

const libDir = "lib"

var separator = []byte("\n\n")

func LoadScript(name string) *gojq.Code {
	return Script(string(mustReadFile(name)))
}

func Script(script string) *gojq.Code {
	buf := bytes.Buffer{}
	readLibFiles(libDir, &buf)
	buf.Write([]byte{'\n'})
	buf.Write([]byte(script))
	return mustCompile(buf.String())
}

func readLibFiles(dir string, buf *bytes.Buffer) {
	entries, err := scripts.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			readLibFiles(dir+"/"+entry.Name(), buf)
		} else if strings.HasSuffix(entry.Name(), ".jq") {
			buf.Write(mustReadFile(dir + "/" + entry.Name()))
			buf.Write(separator)
		}
	}
}

func mustReadFile(name string) []byte {
	file, err := scripts.Open(name)
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return data
}

func mustCompile(query string) *gojq.Code {
	if debug {
		fmt.Fprintln(os.Stderr, "=> compiling jq script \n:"+query)
	}
	parsed, err := gojq.Parse(query)
	if err != nil {
		panic(err)
	}
	compiled, err := gojq.Compile(parsed)
	if err != nil {
		panic(err)
	}
	return compiled
}

type ParsedHook func(any) error

func Run(code *gojq.Code, hook ParsedHook, data []byte) ([]byte, error) {
	var m any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	iter := code.Run(m)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, err
		}
		if hook != nil {
			if err := hook(v); err != nil {
				return nil, err
			}
		}
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		buf.Write(data)
		buf.Write([]byte{'\n'})
	}
	result := buf.Bytes()
	// remove last \n
	if len(data) > 0 && data[len(data)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result, nil
}
