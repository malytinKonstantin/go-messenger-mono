package cassandra

import (
	"github.com/gocql/gocql"
)

// iteratorImpl реализует интерфейс Iterator.
type iteratorImpl struct {
	iter *gocql.Iter
}

func (i *iteratorImpl) Scan(dest ...interface{}) bool {
	return i.iter.Scan(dest...)
}

func (i *iteratorImpl) Close() error {
	if err := i.iter.Close(); err != nil {
		return err
	}
	return nil
}

func (i *iteratorImpl) PageState() []byte {
	return i.iter.PageState()
}

func (i *iteratorImpl) WillSwitchPage() bool {
	return i.iter.WillSwitchPage()
}

func (i *iteratorImpl) Host() *gocql.HostInfo {
	return i.iter.Host()
}
