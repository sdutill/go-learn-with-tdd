package maps

import "testing"

func TestDelete(t *testing.T) {
	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dict := Dictionary{word: definition}

		err := dict.Delete(word)

		assertError(t, err, nil)

		_, err = dict.Search(word)

		assertError(t, err, ErrNotFound)
	})
	t.Run("non-existing word", func(t *testing.T) {
		word := "test"
		dict := Dictionary{}

		err := dict.Delete(word)

		assertError(t, err, ErrWordDoesNotExist)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dict := Dictionary{word: definition}
		newDefinition := "this is an updated test"

		err := dict.Update(word, newDefinition)

		assertError(t, err, nil)
		assertDefinition(t, dict, word, newDefinition)
	})
	t.Run("new word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dict := Dictionary{}

		err := dict.Update(word, definition)

		assertError(t, err, ErrWordDoesNotExist)
		assertDefinition(t, dict, word, definition)
	})
}

func TestAdd(t *testing.T) {
	t.Run("new word", func(t *testing.T) {
		dict := Dictionary{}

		word := "test"
		definition := "this is just a test"
		err := dict.Add(word, definition)

		assertError(t, err, nil)
		assertDefinition(t, dict, "test", "this is just a test")
	})

	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"

		dict := Dictionary{word: definition}

		err := dict.Add(word, definition)

		assertError(t, err, ErrWordExists)
		assertDefinition(t, dict, word, definition)
	})

}

func assertDefinition(t testing.TB, dict Dictionary, word, definition string) {
	t.Helper()

	got, err := dict.Search(word)
	if err != nil {
		t.Fatal("should find added word:", err)
	}
	assertStrings(t, got, definition)
}

func TestSearch(t *testing.T) {
	dict := Dictionary{"test": "this is just a test"}

	t.Run("known word", func(t *testing.T) {
		got, _ := dict.Search("test")
		want := "this is just a test"

		assertStrings(t, got, want)
	})
	t.Run("unknown word", func(t *testing.T) {
		_, got := dict.Search("unknown")

		if got == nil {
			t.Fatal("expeccted to get an error.")
		}

		assertError(t, got, ErrNotFound)
	})
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
