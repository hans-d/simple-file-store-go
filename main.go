package filestore

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

var (
	ErrMissingKey    = errors.New("missing key - unable to save data")
	ErrMissingParent = errors.New("missing parent - no place to save data")
)

// Logger is a generic logger
type Logger interface {
	Fatal(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
}

// Driver holds the config and interacts with the underlying file store.
type Driver struct {
	mutex     sync.Mutex
	mutexes   map[string]*sync.Mutex
	dir       string // where files are stored
	log       Logger
	marshaler Marshaler // file format
}

// Options for optional config
type Options struct {
	Logger
	Marshaler
}

// New creates a new driver at the given location, and returns a *Driver
// for further interaction. By default will use teh JSONMarshaler.
func New(dir string, options *Options) (*Driver, error) {

	if options == nil {
		options = &Options{}
	}
	if options.Logger == nil {
		options.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}
	if options.Marshaler == nil {
		options.Marshaler = &JSONMarshaler{}
	}
	dir = filepath.Clean(dir)

	driver := Driver{
		dir:       dir,
		mutexes:   make(map[string]*sync.Mutex),
		log:       options.Logger,
		marshaler: options.Marshaler,
	}
	if _, err := os.Stat(dir); err == nil {
		options.Logger.Debug("Using '%s' (folder already exists)\n", dir)
		return &driver, nil
	}

	options.Logger.Debug("Creating storage folder at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

// Write the value [v] to the file [key] under [parent], using the
// marshaler. A lock is held on [parent].
func (d *Driver) Write(parent, key string, v interface{}) error {

	if parent == "" {
		return ErrMissingParent
	}
	if key == "" {
		return ErrMissingKey
	}

	mutex := d.getOrCreateMutex(parent)
	mutex.Lock()
	defer mutex.Unlock()

	return d.writeFile(parent, key, v)
}

func (d *Driver) writeFile(parent, key string, v interface{}) error {

	dir := filepath.Join(d.dir, parent)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := d.marshaler.Marshal(v)
	if err != nil {
		return err
	}

	file := filepath.Join(dir, key+d.marshaler.GetFileExtension())
	tmpFile := file + ".tmp"

	if err := ioutil.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmpFile, file)
}

// Read the content from [key] under [parent] into [v].
func (d *Driver) Read(parent, key string, v interface{}) error {

	if parent == "" {
		return ErrMissingParent
	}
	if key == "" {
		return ErrMissingKey
	}

	return d.readFile(parent, key, v)
}

func (d *Driver) readFile(parent, key string, v interface{}) error {

	file := filepath.Join(d.dir, parent, key+d.marshaler.GetFileExtension())

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return d.marshaler.Unmarshal(data, v)
}

func (d *Driver) getOrCreateMutex(parent string) *sync.Mutex {

	d.mutex.Lock()
	defer d.mutex.Unlock()

	m, ok := d.mutexes[parent]
	if !ok {
		m = &sync.Mutex{}
		d.mutexes[parent] = m
	}

	return m
}
