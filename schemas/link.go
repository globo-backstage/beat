package schemas

import (
	"fmt"
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

type Links []*Link

func (l *Links) ApplyBaseUrl(baseUrl string) {
	for _, link := range *l {
		link.Href = fmt.Sprintf("%s%s", baseUrl, link.Href)
	}
}

func BuildDefaultLinks(collectionName string) Links {
	collectionUrl := fmt.Sprintf("/%s", collectionName)
	itemUrl := fmt.Sprintf("/%s/{id}", collectionName)

	return Links{
		&Link{Rel: "self", Href: itemUrl},
		&Link{Rel: "item", Href: itemUrl},
		&Link{Rel: "create", Method: "POST", Href: collectionUrl},
		&Link{Rel: "update", Method: "PUT", Href: itemUrl},
		&Link{Rel: "delete", Method: "DELETE", Href: itemUrl},
		&Link{Rel: "parent", Href: collectionUrl},
	}
}