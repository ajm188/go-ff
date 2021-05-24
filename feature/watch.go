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
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	base := filepath.Base(path)

	log.Printf("[watch] adding %s", dir)
	if err := watcher.Add(dir); err != nil {
		watcher.Close()
		return err
	}

	go func() {
		defer watcher.Close()

		for {
			select {
			case <-ctx.Done():
				log.Printf("[watch] context finished: %v", ctx.Err())
				return
			case event, ok := <-watcher.Events:
				if !ok {
					log.Print("[watch] events channel closed, stopping watch")
					return
				}

				if filepath.Base(event.Name) != base {
					continue
				}

				if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
					log.Printf("[watch] config touched but not written (op = %v), ignoring", event.Op)
					continue
				}

				log.Print("[watch] detected config change. reloading ...")

				if err := InitFromFile(event.Name); err != nil { // TODO: add noglog
					// Purely for logging considerations, distinguish between
					// errors that already contain the pathname vs those that
					// don't.
					if errors.Is(err, ErrInvalidFeature) || errors.Is(err, ErrUnknownFeatureType) {
						log.Printf("error reloading config from %s: %v", event.Name, err)
					} else {
						log.Printf("error reloading config: %v", err)
					}
				}

				log.Print("[watch] reloaded config")
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Print("[watch] errors channel closed, stopping watch")
					return
				}

				// TODO: add noglog
				log.Printf("error watching %s: %v", dir, err)
			}
		}
	}()

	return nil
}
