from cliquet import resource


@resource.register(name='item-schema',
                   collection_path='/api/item-schema',
                   record_path='/api/item-schema/{{collectionName}}')
class ItemSchema(resource.UserResource):
    # No schema yet.
    pass
