package config

import "strings"

type Path string

func (s *Path) UnmarshalYAML(unmarshaller func(any) error) error {
	var v string
	if err := unmarshaller(&v); err != nil {
		return err
	}
	v = strings.TrimSpace(v)

	for strings.Contains(v, "//") {
		v = strings.ReplaceAll(v, "//", "/")
	}

	v = strings.TrimSuffix(v, "/")

	*s = Path(v)
	return nil
}
