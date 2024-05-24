// Taken from https://samrapdev.com/capturing-sensitive-input-with-editor-in-golang-from-the-cli/

package editor

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

var ErrNoChange = errors.New("no changes made to file")

const DefaultEditor = "vim"

// OpenFileInEditor opens the file at the given path for editing with the user's
// default editor. The default editor is determined via the `EDITOR` env
// variable. If not set, the default editor (vim) is used.
func OpenFileInEditor(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}

	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// EditData writes the given data to a temporary file, then calls
// `OpenFileInEditor` to edit the data. It returns the edited data and any
// errors encountered while editing the data. If the given data was not
// modified, `ErrNoChange` is returned.
func EditData(data []byte) ([]byte, error) {
	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		return nil, err
	}

	defer os.Remove(file.Name())

	if _, err := io.Copy(file, bytes.NewReader(data)); err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	if err = OpenFileInEditor(file.Name()); err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return nil, err
	}

	if !modified(data, bytes) {
		return data, ErrNoChange
	}

	return bytes, nil
}

func modified(old, new []byte) bool {
	hash := md5.New()

	io.Copy(hash, bytes.NewReader(old))

	oldHash := hex.EncodeToString(hash.Sum(nil)[:16])

	hash.Reset()

	io.Copy(hash, bytes.NewReader(new))

	newHash := hex.EncodeToString(hash.Sum(nil)[:16])

	return oldHash != newHash
}
