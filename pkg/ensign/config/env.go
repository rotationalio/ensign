package config

import (
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

// Ensure that the environment is only loaded once.
var (
	loadenv sync.Once
	envvars map[string]struct{}
)

// Merges the values from envconf into conf if the conf value is zero-valued or if an
// environment variable exists for the envconf field. Because envconf does not export a
// method that will allow us to identify the environment variable we fetch candidate
// environment variables with the Ensign prefix and inspect the struct tags to do a best
// effort guess if the environment variable is there. This is likely error prone and to
// do this correctly we need to port the envconfig code to this library.
func mergenv(conf, envconf interface{}) error {
	// cs refers to "conf spec" and es to "envconf spec" to get reflection started.
	cs := reflect.ValueOf(conf)
	es := reflect.ValueOf(envconf)

	if cs.Kind() != reflect.Ptr || es.Kind() != reflect.Ptr {
		return envconfig.ErrInvalidSpecification
	}

	return merge(prefix, cs.Elem(), es.Elem())
}

// Recursive merging function that expects cs and es to be a config struct of the same
// type. The values from es are merged into cs if the cs value is zero-valued or an
// environment variable candidate exists for the es field.
func merge(prefix string, cs, es reflect.Value) error {
	if cs.Kind() != reflect.Struct || es.Kind() != reflect.Struct {
		return envconfig.ErrInvalidSpecification
	}

	if cs.Type() != es.Type() {
		return envconfig.ErrInvalidSpecification
	}

	specType := es.Type()

fields:
	for i := 0; i < cs.NumField(); i++ {
		// cs refers to the conf field and ef to the envconf field.
		cf := cs.Field(i)
		ef := es.Field(i)
		ftype := specType.Field(i)

		if !cf.CanSet() {
			continue fields
		}

		for cf.Kind() == reflect.Ptr {
			if ef.IsNil() {
				// There is nothing in the env conf field so stop processing
				continue fields
			}

			if cf.IsNil() {
				if cf.Type().Elem().Kind() != reflect.Struct {
					// nil pointer to a non-struct; leave it alone
					break
				}

				// nil pointer to a struct; create a zero-valued instance
				cf.Set(reflect.New(cf.Type().Elem()))
			}

			cf = cf.Elem()
			ef = ef.Elem()
		}

		// If the field is a struct, recursively merge
		if cf.Kind() == reflect.Struct {
			if err := merge(prefix+ftype.Name, cf, ef); err != nil {
				return err
			}
			continue fields
		}

		// Attempt to determine what the key might be
		key := ftype.Tag.Get("envconfig")
		if key == "" {
			key = prefix + ftype.Name
		}

		// Otherwise set the cf field from the ef field
		// TODO: how to perform the environment variable check?
		if cf.IsZero() || checkEnv(key) {
			cf.Set(ef)
		}
	}

	return nil
}

// Returns true if the given key is in the environment. All keys are normalized by
// removing underscores and making them uppercase to help with the search.
func checkEnv(key string) bool {
	// Ensure the env is loaded once before we perform our check
	loadenv.Do(func() {
		envvars = make(map[string]struct{})
		for _, pair := range os.Environ() {
			key := strings.Split(pair, "=")[0]
			envvars[strings.ToUpper(strings.Replace(key, "_", "", -1))] = struct{}{}
		}
	})

	key = strings.ToUpper(strings.Replace(key, "_", "", -1))
	_, ok := envvars[key]
	return ok
}

// Reset the environment used by checkEnv (primarily for testing).
func ResetLocalEnviron() {
	loadenv = sync.Once{}
	envvars = nil
}
