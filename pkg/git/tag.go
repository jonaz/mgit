package git

import (
	"strings"

	"github.com/jonaz/mgit/pkg/utils"
)

type Tag struct {
	Hash string
	Tag  string
}

type Tags map[string]Tag

func (t Tags) ByHash(hash string) *Tag {
	for _, tag := range t {
		if tag.Hash == hash {
			t := tag
			return &t
		}
	}
	return nil
}

func (t Tags) ByTag(tag string) *Tag {
	if tag, ok := t[tag]; ok {
		t := tag
		return &t
	}
	return nil
}

// Tags lists all available tags with cheir corresponding git hash.
func ListTags() (Tags, error) {
	tags := make(Tags)
	output, err := utils.Run("git", "for-each-ref", "--format=%(if)%(*objectname)%(then)%(*objectname)%(else)%(objectname)%(end) %(refname)", "refs/tags")
	if err != nil {
		return nil, err
	}

	for _, v := range strings.Split(output, "\n") {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		tmp := strings.Split(v, " ")
		tag := strings.Split(tmp[1], "/")[2]

		tags[tag] = Tag{
			Hash: tmp[0],
			Tag:  tag,
		}
	}
	return tags, nil
}
