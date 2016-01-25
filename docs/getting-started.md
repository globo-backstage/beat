# Installing

### Requirements

- Go 1.5+
- MongoDB 3+

### Download and install the devolpement version

Ensure if your GOPATH environment variable is setted, see more in: https://golang.org/doc/code.html#GOPATH

```bash
go get "github.com/backstage/beat/beat"
cd $GOPATH/src/github.com/backstage/beat
make setup
```

### Run the devolpement version

```
make run
```

# Using (with Curl)

### Create a new collection

To dynamically define a new collection just create a new instance of the ItemSchema. Doing this via the REST interface is as simples as POSTing a valid JSON Schema, as follows:

###### schema.json
```json
{
  "collectionName": "people",
  "globalCollectionName": true,
  "type": "object",
  "title": "Person",
  "collectionTitle": "People",
  "properties": {
    "name": {
      "type": "string"
    },
    "email": {
      "type": "string",
      "format": "email"
    }
  }
}
```

Create a Person model from a JSON Schema
```
curl -i -XPOST -H "Content-Type: application/json" http://beat-service-example.org/api/item-schemas -T schema.json
```

The REstful API will then be available at http://beat-service-example.org/api/people.

#### Default links

Each Item schema have a default set of links which correspond to the basic CRUD operations supported by Backstage-Beat.

For example:

```
$ curl http://beat-service-example.org/api/item-schemas/people
```

returns
```json
{
  "$schema": "http://json-schema.org/draft-04/hyper-schema#",
  "collectionName": "people",
  ...
  "links": [
    {
      "rel": "self",
      "href": "http://beat-service-example.org/api/people/{id}"
    },
    {
      "rel": "item",
      "href": "http://beat-service-example.org/api/people/{id}"
    },
    {
      "rel": "create",
      "href": "http://beat-service-example.org/api/people",
      "method": "POST",
      "schema": {
        "$ref": "http://beat-service-example.org/api/item-schemas/people"
      }
    },
    {
      "rel": "update",
      "href": "http://beat-service-example.org/api/people/{id}",
      "method": "PUT"
    },
    {
      "rel": "delete",
      "href": "http://beat-service-example.org/api/people/{id}",
      "method": "DELETE"
    },
    {
      "rel": "parent",
      "href": "http://beat-service-example.org/api/people"
    }
  ]
}
```

#### Including custom links in an item schema

It is possible to include custom links in an item schema. To do so, just include them in the links property:

```json
{
  "type": "object",
  ...
  "properties": {
    ...
  },
  "links": [
    {
      "rel": "my-custom-item-schema-link",
      "href": "http://example.org/my/custom/item-schema-link"
    }
  ]
}
```
