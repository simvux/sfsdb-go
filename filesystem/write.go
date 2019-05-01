package filesystem

import (
	"github.com/vmihailenco/msgpack"
	"os"
)

func Save(path Filepath, data interface{}) error {
	f, err := os.Create(path.Unwrap())
	if err != nil {
		return err
	}
	enc := msgpack.NewEncoder(f)
	return enc.Encode(data)
}

func Delete(path Filepath) error {
	return os.Remove(path.Unwrap())
}
