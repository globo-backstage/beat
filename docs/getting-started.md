# Installing

## Requirements

- Go 1.5+
- MongoDB 3+

## Download and install the development version

> Ensure that you have your `GOPATH` environment variable properly configured. Check the [Go docs](https://golang.org/doc/code.html#GOPATH) to see how to id:

```bash
go get "github.com/backstage/beat/beat"
cd $GOPATH/src/github.com/backstage/beat
make setup
```

## Running locally

```
make run
```

## Using (with `curl`)

### Create a new collection

To dynamically define a new collection, just create a new instance of the `ItemSchema`. You can do this using the REST interface to `POST` a valid JSON Schema. First define your schema as below:

##### `schema.json`

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

Then you can create a `Person` collection by POSTing the JSON Schema above:

```bash
curl -i -XPOST -H "Content-Type: application/json" http://beat-service-example.org/api/item-schemas -T schema.json
```

That is it. The RESTful API will then be available at http://beat-service-example.org/api/people.

#### Default links

Each Item schema have a default set of links which correspond to the basic CRUD operations supported by Backstage Beat. For example:

```bash
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

#### Including custom links in an Item Schema

It is possible to include custom links in an Item Schema. To do so, just include them in the links property of your JSON:

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
