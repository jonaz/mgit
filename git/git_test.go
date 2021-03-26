package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByHash(t *testing.T) {
	tags := make(Tags)

	tags["v1.0.2"] = Tag{Hash: "926371aa82ea8ddcb5440dbe215cc9123fb4c333", Tag: "v1.0.2"}
	tags["v1.0.3"] = Tag{Hash: "9a86ed6561d294587d864bda376aa470a0dcc9c6", Tag: "v1.0.3"}
	tags["v1.0.4"] = Tag{Hash: "81da241302ecf8e0efa73b06026846b918bc09d6", Tag: "v1.0.4"}

	tag := tags.ByHash("926371aa82ea8ddcb5440dbe215cc9123fb4c333")

	t.Log(tag)
	assert.Equal(t, tag.Tag, "v1.0.2")
}

func TestByTag(t *testing.T) {
	tags := make(Tags)

	tags["v1.0.2"] = Tag{Hash: "926371aa82ea8ddcb5440dbe215cc9123fb4c333", Tag: "v1.0.2"}
	tags["v1.0.3"] = Tag{Hash: "9a86ed6561d294587d864bda376aa470a0dcc9c6", Tag: "v1.0.3"}
	tags["v1.0.4"] = Tag{Hash: "81da241302ecf8e0efa73b06026846b918bc09d6", Tag: "v1.0.4"}

	tag := tags.ByTag("v1.0.3")

	t.Log(tag)
	assert.Equal(t, tag.Hash, "9a86ed6561d294587d864bda376aa470a0dcc9c6")
}
