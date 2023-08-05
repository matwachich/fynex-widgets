package wx

import (
	"embed"
	"io/fs"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"
)

type EmbedRepository struct {
	scheme string
	fs     embed.FS
}

func RegisterEmbedRepository(scheme string, fs embed.FS) {
	repository.Register(scheme, &EmbedRepository{scheme: scheme, fs: fs})
}

type readCloser struct {
	fs.File
	u fyne.URI
}

func (rc *readCloser) URI() fyne.URI { return rc.u }

func (e *EmbedRepository) checkScheme(u fyne.URI) bool {
	return u.Scheme() == e.scheme
}

func (e *EmbedRepository) openURI(u fyne.URI) (rc *readCloser, err error) {
	var f fs.File
	if f, err = e.fs.Open(strings.TrimPrefix(u.Path(), "/")); err != nil {
		return
	}
	return &readCloser{File: f, u: u}, nil
}

// Repository

// Exists will be used to implement calls to storage.Exists() for the
// registered scheme of this repository.
//
// Since: 2.0
func (e *EmbedRepository) Exists(u fyne.URI) (bool, error) {
	if !e.checkScheme(u) {
		return false, repository.ErrOperationNotSupported
	}

	rc, err := e.openURI(u)
	if err != nil {
		return false, err
	}
	defer rc.Close()

	return true, nil
}

// Reader will be used to implement calls to storage.Reader()
// for the registered scheme of this repository.
//
// Since: 2.0
func (e *EmbedRepository) Reader(u fyne.URI) (fyne.URIReadCloser, error) {
	if !e.checkScheme(u) {
		return nil, repository.ErrOperationNotSupported
	}

	return e.openURI(u)
}

// CanRead will be used to implement calls to storage.CanRead() for the
// registered scheme of this repository.
//
// Since: 2.0
func (e *EmbedRepository) CanRead(u fyne.URI) (bool, error) {
	if !e.checkScheme(u) {
		return false, repository.ErrOperationNotSupported
	}

	rc, err := e.openURI(u)
	if err != nil {
		return false, err
	}
	defer rc.Close()

	stat, err := rc.Stat()
	if err != nil {
		return false, err
	}

	return !stat.IsDir(), nil
}

// Destroy is called when the repository is un-registered from a given
// URI scheme.
//
// The string parameter will be the URI scheme that the repository was
// registered for. This may be useful for repositories that need to
// handle more than one URI scheme internally.
//
// Since: 2.0
func (e *EmbedRepository) Destroy(string) {}

// ListableRepository

// CanList will be used to implement calls to storage.Listable() for
// the registered scheme of this repository.
//
// Since: 2.0
func (e *EmbedRepository) CanList(u fyne.URI) (bool, error) {
	if !e.checkScheme(u) {
		return false, repository.ErrOperationNotSupported
	}

	rc, err := e.openURI(u)
	if err != nil {
		return false, err
	}
	defer rc.Close()

	stat, err := rc.Stat()
	if err != nil {
		return false, err
	}

	return stat.IsDir(), nil
}

// List will be used to implement calls to storage.List() for the
// registered scheme of this repository.
//
// Since: 2.0
func (e *EmbedRepository) List(u fyne.URI) ([]fyne.URI, error) {
	if !e.checkScheme(u) {
		return nil, repository.ErrOperationNotSupported
	}

	list, err := e.fs.ReadDir(strings.TrimLeft(u.Path(), "/"))
	if err != nil {
		return nil, err
	}

	ret := make([]fyne.URI, 0, len(list))
	for _, elem := range list {
		u, _ := storage.ParseURI(u.String() + "/" + elem.Name())
		ret = append(ret, u)
	}
	return ret, nil
}

// CreateListable will be used to implement calls to
// storage.CreateListable() for the registered scheme of this
// repository.
//
// Since: 2.0
func (e *EmbedRepository) CreateListable(u fyne.URI) error { return nil }
