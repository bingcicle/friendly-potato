package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type anymap map[string]any

func readJSONL(path string) ([]anymap, error) {
	f, err := os.Open(path)
	if err != nil { return nil, err }
	defer f.Close()
	var out []anymap
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" { continue }
		var m anymap
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		out = append(out, m)
	}
	return out, sc.Err()
}

func writeJSONL(items []anymap) error {
	w := bufio.NewWriter(os.Stdout)
	for _, m := range items {
		b, _ := json.Marshal(m)
		fmt.Fprintln(w, string(b))
	}
	return w.Flush()
}

func get(m anymap, key string) string {
	v, ok := m[key]
	if !ok || v == nil { return "" }
	return fmt.Sprintf("%v", v)
}

func main() {
	mode := flag.String("mode", "merge", "merge|filter")
	inputs := flag.String("in", "", "comma-separated JSONL files")
	field := flag.String("field", "", "field name for filter")
	eq := flag.String("eq", "", "value equals (string compare)")
	rex := flag.String("rex", "", "regex match")
	flag.Parse()

	if *inputs == "" { log.Fatal("--in required") }

	var all []anymap
	for _, p := range strings.Split(*inputs, ",") {
		items, err := readJSONL(strings.TrimSpace(p))
		if err != nil { log.Fatal(err) }
		all = append(all, items...)
	}

	switch *mode {
	case "merge":
		if err := writeJSONL(all); err != nil { log.Fatal(err) }
	case "filter":
		if *field == "" { log.Fatal("--field required for filter") }
		var out []anymap
		var r *regexp.Regexp
		if *rex != "" {
			var err error
			r, err = regexp.Compile(*rex)
			if err != nil { log.Fatal(err) }
		}
		for _, m := range all {
			val := get(m, *field)
			ok := true
			if *eq != "" { ok = ok && (val == *eq) }
			if r != nil { ok = ok && r.MatchString(val) }
			if ok { out = append(out, m) }
		}
		if err := writeJSONL(out); err != nil { log.Fatal(err) }
	default:
		log.Fatal("unknown --mode")
	}
}
