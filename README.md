# fugozi

### what?
Terribly named in-memory document database (read: key:value), pretty much expecting json,
though really any text will work.
The general idea here is that you can one or many `bucket` with (n)one or many
`docs` under each bucket. I.e a bucket is just a collection of docs.
A doc is really a key:value pair, the key being the `docid`, the value being basically
text, most likely in json format.

The database itself is accessible via an HTTP interface.
```
GET /status                   : prints a crappy status page in json Format
POST /bucket/$bucketid        : creates a bucket with id of $bucketid
GET /bucket/$bucketid         : prints a bucket name/id in json format
POST /bucket/$bucketid/$docid : creates a document $docid under bucket $bucketid.
                                we are expecting a request body to be stored as the
                                $docid's doc.
GET /bucket/$bucketid/$docid  : sends back whatever you stored at $docid under $bucketid
```

404s and 500s abound when you do something wrong.

Yes, the underlying maps / collections are concurrent.
