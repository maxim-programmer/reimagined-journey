package chunker

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
    t.Run("empty string returns nil", func(t *testing.T) {
        result := Split("")
        assert.Nil(t, result)
    })

    t.Run("short text returns one chunk", func(t *testing.T) {
        text := "Hello, world!"
        result := Split(text)
        assert.Len(t, result, 1)
        assert.Equal(t, text, result[0])
    })

    t.Run("long text splits into multiple chunks", func(t *testing.T) {
        text := make([]byte, 1500)
        for i := range text {
            text[i] = 'a'
        }
        result := Split(string(text))
        assert.Greater(t, len(result), 1)
    })

    t.Run("cyrillic text preserves characters", func(t *testing.T) {
        text := "Привет мир! " + string(make([]byte, 500))
        result := Split(text)
        for _, chunk := range result {
            assert.NotContains(t, chunk, "�")
        }
    })
}

func TestSplitPages(t *testing.T) {
    t.Skip("requires model package import")
}
