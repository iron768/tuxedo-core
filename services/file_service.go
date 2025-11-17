package services

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
    watcher *fsnotify.Watcher
    path    string
}

func NewFileWatcher(path string) (*FileWatcher, error) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return nil, err
    }
    
    if err := watcher.Add(path); err != nil {
        return nil, err
    }
    
    return &FileWatcher{watcher: watcher, path: path}, nil
}

func (fw *FileWatcher) Watch(callback func(event fsnotify.Event)) {
    go func() {
        for {
            select {
            case event := <-fw.watcher.Events:
                callback(event)
            case err := <-fw.watcher.Errors:
                log.Println("Watcher error:", err)
            }
        }
    }()
}