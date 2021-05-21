package feature

import (
	"context"
	"errors"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Watch watches the given path for changes and reloads the global feature
// config.
//
// In reality, Watch watches the directory of the given path, rather than the
// actual filepath, to handle cases where the file is deleted. Filesystem events
// in that directory unrelated to the particular filepath are ignored.
//
// When the config path is modified, Watch uses InitFromFile to read the file,
// ensure it is non-empty, unmarshal it from json, and validate the feature
// specs before swapping in the config.
//
// The watch continues until the watcher closes either the Events or Errors
// channels, or until the context is cancelled or expired.
func Watch(ctx context.Context, path string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer watcher.Close()

	dir := filepath.Dir(path)
	base := filepath.Base(path)

	if err := watcher.Add(dir); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if filepath.Base(event.Name) != base {
					continue
				}

				if err := InitFromFile(event.Name); err != nil { // TODO: add noglog
					// Purely for logging considerations, distinguish between
					// errors that already contain the pathname vs those that
					// don't.
					if errors.Is(err, ErrInvalidFeature) || errors.Is(err, ErrUnknownFeatureType) {
						log.Printf("error reloading config from %s: %v", event.Name, err)
					} else {
						log.Printf("error reloading config %v", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				// TODO: add noglog
				log.Printf("error watching %s: %v", dir, err)
			}
		}
	}()

	return nil
}
