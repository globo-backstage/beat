package schemas

import (
	"fmt"
	"net/url"
)

type Link struct {
	Rel          string                 `json:"rel" bson:"rel"`
	Href         string                 `json:"href" bson:"href"`
	Title        string                 `json:"title,omitempty" bson:"title,omitempty"`
	TargetSchema map[string]interface{} `json:"targetSchema,omitempty" bson:"targetSchema,omitempty"`
	MediaType    string                 `json:"mediaType,omitempty" bson:"mediaType,omitempty"`
	Method       string                 `json:"method,omitempty" bson:"method,omitempty"`
	EncType      string                 `json:"encType,omitempty" bson:"encType,omitempty"`
	Schema       map[string]interface{} `json:"schema,omitempty" bson:"schema,omitempty"`
}

var (
	DefaultLinkRels = []string{"self", "item", "create", "update", "delete", "parent"}
)

type Links []*Link

func (l Links) ApplyBaseUrl(baseUrl string) {
	for _, link := range l {
		if isRelativeLink(link.Href) {
			link.Href = fmt.Sprintf("%s%s", baseUrl, link.Href)
		}

		if ref, ok := link.Schema["$ref"].(string); ok && isRelativeLink(ref) {
			link.Schema["$ref"] = fmt.Sprintf("%s%s", baseUrl, ref)
		}

		if ref, ok := link.TargetSchema["$ref"].(string); ok && isRelativeLink(ref) {
			link.TargetSchema["$ref"] = fmt.Sprintf("%s%s", baseUrl, ref)
		}
	}
}

// ConcatenateLinks generate new links with merge with tailLinks
func (l Links) ConcatenateLinks(tailLinks *Links) *Links {
	currentSize := len(l)
	expandSize := len(*tailLinks)

	newLinks := make(Links, currentSize+expandSize)
	copy(newLinks, l)

	for i, link := range *tailLinks {
		newLinks[currentSize+i] = link
	}

	return &newLinks
}

// DiscardDefaultLinks remove all default links to store only custom links
func (l Links) DiscardDefaultLinks() *Links {
	newLinks := make(Links, 0, len(l))
	for _, link := range l {
		if !isDefaultRel(link.Rel) {
			newLinks = append(newLinks, link)
		}
	}
	return &newLinks
}

func isRelativeLink(link string) bool {
	url, err := url.Parse(link)

	if err != nil {
		return false
	}

	return url.Host == "" && url.Scheme == "" && !isUriTemplate(link)
}

func isUriTemplate(link string) bool {
	return len(link) > 0 && link[0] == '{'
}

func isDefaultRel(linkRel string) bool {
	for _, defaultLinkRel := range DefaultLinkRels {
		if defaultLinkRel == linkRel {
			return true
		}
	}
	return false
}
