package main

type Dictionary map[string]string
type DictionaryErr string

const (
	ErrNotFound         = DictionaryErr("could not find the word you were looking for")
	ErrWordExists       = DictionaryErr("cannot add word because it already exists")
	ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

func (e DictionaryErr) Error() string {
	return string(e)
}

func (d Dictionary) Search(word string) (string, error) {
	definition, ok := d[word]
	if !ok {
		return "", ErrNotFound
	}

	return definition, nil
}

func (d Dictionary) Add(word, definition string) error {
	_, err := d.Search(word)

	switch err { // Switch is based on the result of the error returned by Search
	case ErrNotFound: // word not found in dict, so we can add it
		d[word] = definition
	case nil: // word already exists in dict, so we return an error
		return ErrWordExists
	default:
		return err
	}
	return nil
}

func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)

	switch err {
	case ErrNotFound: // word not found in dict, so we cannot update it
		return ErrWordDoesNotExist
	case nil: // word exists in dict, so we can update it
		d[word] = definition
	default:
		return err
	}

	return nil
}

func (d Dictionary) Delete(word string) error {
	delete(d, word)
	_, err := d.Search(word)
	switch err {
	case ErrNotFound: // word not found in dict, so we cannot delete it
		return nil
	case nil: // word exists in dict, so we can delete it
		delete(d, word)
	default:
		return err
	}

	return nil
}
