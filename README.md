VecDB
======
> a very simple vector embedding database, 
> you can say that it is a hash-table that let you find items similar to the item you're searching for.

Why!
====
> I'm a databases enthusiast, and this is a for fun and learning project that could be used in production ;).
> 
> **P.S**: I like to re-invent the wheel in my free time, because it is my free time!

Data Model
==========
> I'm using the `{key => value}` model,
> - `key` should be a unique value that represents the item.
> - `value` should be the vector itself (List of Floats).

Configurations
==============
> by default `vecdb` searches for `config.yml` in the current working directory.
> but you can override it using the `--config /path/to/config.yml` flag by providing your own custom file path.

```yaml
# http server related configs
server:
  # the address to listen on in the form of '[host]:port'
  listen: "0.0.0.0:3000"

# storage related configs
store:
  # the driver you want to use
  # currently vecdb supports "bolt" which is based on boltdb the in process embedded the database
  driver: "bolt"
  # the arguments required by the driver
  # for bolt, it requires a key called `database` points to the path you want to store the data in.
  args:
    database: "./vec.db"

# embeddings related configs
embedder:
  # whether to enable the embedder and all endpoints using it or not
  enabled: true
  # the driver you want to use, currently vecdb supports gemini
  driver: gemini
  # the arguments required by the driver
  # currently gemini driver requires `api_key` and `text_embedding_model`
  args:
    # by default vecdb will replace anything between ${..} with the actual value from the ENV var
    api_key: "${GEMINI_API_KEY}"
    text_embedding_model: "text-embedding-004"
```

Components
===========
- Raw Vectors Layer (low-level)
  - send [VectorWriteRequest](#VectorWriteRequest) to `POST /v1/vectors/write` when you have a vector and want to store it somewhere.
  - send [VectorSearchRequest](#VectorSearchRequest) to `POST /v1/vectors/search` when you have a vector and want to list all similar vectors' keys/ids ordered by cosine similarity in descending order.
- Embedding Layer (optional)
  - send [TextEmbeddingWriteRequest](#TextEmbeddingWriteRequest) to `POST /v1/embeddings/text/write` when you have a text and want `vecdb` to build and store the vector for you using the configured embedder (gemini for now).
  - send [TextEmbeddingSearchRequest](#TextEmbeddingSearchRequest) to `POST /v1/embeddings/text/search` when you have a text and want `vecdb` to build a vector and search for similar vectors' keys for you ordered by cosine similarity in descending order.

Requests
========

### VectorWriteRequest
```json5
{
  "bucket": "BUCKET_NAME", // consider it a collection or a table
  "key": "product-id-1", // should be unique and represents a valid value in your main data store (example: the row id in your mysql/postgres ... etc)
  "vector": [1.929292, 0.3848484, -1.9383838383, ... ] // the vector you want to store 
}
```

### VectorSearchRequest
```json5
{
  "bucket": "BUCKET_NAME", // consider it a collection or a table
  "vector": [1.929292, 0.3848484, -1.9383838383, ... ], // you will get a list ordered by cosine-similarity in descending order
  "min_cosine_similarity": 0.0, // the more you increase, the fewer data you will get
  "max_result_count": 10 // max vectors to return (vecdb will first order by cosine similarity then apply the limit)
}
```

### TextEmbeddingWriteRequest
> if you set `embedder.enabled` to `true`.

```json5
{
  "bucket": "BUCKET_NAME", // consider it a collection or a table
  "key": "product-id-1", // should be unique and represents a valid value in your main data store (example: the row id in your mysql/postgres ... etc)
  "content": "This is some text representing the product" // this will be converted to a vector using the configured embedder 
}
```

### TextEmbeddingSearchRequest
> if you set `embedder.enabled` to `true`.

```json5
{
  "bucket": "BUCKET_NAME", // consider it a collection or a table
  "content": "A Product Text", // you will get a list ordered by cosine-similarity in descending order
  "min_cosine_similarity": 0.0, // the more you increase, the fewer data you will get
  "max_result_count": 10 // max vectors to return (vecdb will first order by cosine similarity then apply the limit)
}
```

Download/Install
================
- [Binary](https://github.com/alash3al/vecdb/releases)
- [Docker Image](https://github.com/alash3al/vecdb/pkgs/container/vecdb)