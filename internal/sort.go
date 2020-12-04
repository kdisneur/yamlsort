package internal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

// CustomLabel is the key used when the yaml to indent doesn't start with an object key
const CustomLabel string = "YAMLSORT_BEGIN_KEY:"

// ErrNoStartingLabel is an error used when the YAML selection doesn't start with an object
var ErrNoStartingLabel = errors.New("doesn't start with an object")

// SortYAML sorts a complete or partial YAML file
func SortYAML(in io.Reader, out io.Writer, indent int) error {

	incomingYAML, err := ioutil.ReadAll(in)
	if err != nil {
		return fmt.Errorf("can't read input: %v", err)
	}

	var hasNoStartingLabel bool
	rootIndent, err := detectRootIndent(incomingYAML)
	if err != nil {
		if !errors.Is(err, ErrNoStartingLabel) {
			fmt.Fprint(out, string(incomingYAML))
			return fmt.Errorf("can't detect root indentation: %v", err)
		}

		hasNoStartingLabel = true
	}

	if hasNoStartingLabel {
		incomingYAML = append([]byte(CustomLabel+"\n"), incomingYAML...)
	}

	var value map[string]interface{}
	if err := yaml.Unmarshal(incomingYAML, &value); err != nil {
		fmt.Fprint(out, string(incomingYAML))

		return fmt.Errorf("can't decode YAML: %v", err)
	}

	var outgoingYAML bytes.Buffer
	encoder := yaml.NewEncoder(&outgoingYAML)
	encoder.SetIndent(indent)

	if err := encoder.Encode(&value); err != nil {
		fmt.Fprint(out, string(incomingYAML))
		return fmt.Errorf("can't re-encode YAML: %v", err)
	}

	reindentedYAML, err := indentYAML(outgoingYAML.String(), rootIndent, indent, hasNoStartingLabel)
	if err != nil {
		fmt.Fprint(out, string(incomingYAML))
		return fmt.Errorf("can't re-indent YAML: %v", err)
	}

	fmt.Fprint(out, reindentedYAML)
	return nil
}

func indentYAML(yaml string, rootIndent string, indent int, hasNoStartingLabel bool) (string, error) {
	var newYAML strings.Builder

	var count int
	scanner := bufio.NewScanner(strings.NewReader(yaml))
	for scanner.Scan() {
		count++

		if count == 1 && hasNoStartingLabel {
			continue
		}

		line := scanner.Text()
		if hasNoStartingLabel {
			line = line[2:]
		}

		newYAML.WriteString(rootIndent + line + "\n")
	}

	if scanner.Err() != nil {
		return "", fmt.Errorf("can'read line: %v", scanner.Err())
	}

	return newYAML.String(), nil
}

func detectRootIndent(yaml []byte) (string, error) {
	var indent string

	reader := bufio.NewReader(bytes.NewReader(yaml))
	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			return "", fmt.Errorf("can't read character: %v", err)
		}

		if char == ' ' || char == '\t' {
			indent += string(char)
			continue
		}

		if char == '-' {
			return indent, ErrNoStartingLabel
		}

		break
	}

	return indent, nil
}
