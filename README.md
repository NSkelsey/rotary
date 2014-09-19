Rotary
======

A small go server that functions as a hash to file blob store. Storage is done in sqlite!

__Note__! There is a 10MB file size limit


##Routes

`POST /upload`

Upload some object via a http post. 
- Requires a Content-Type header providing the correct MIME type for the item.
- Requires a Content-Length header set to the length of the request's body.
- Returns 201 on success and greater than 400 on failure.

`GET /<base64_urlencode>` 

Returns the raw file with content-type set by the header on.

`GET /api/<base64_urlencode>`

Functions like the bare / GET except the server returns a json object that looks like:
```json
{
   "Raw" : "d29vb29vb29vb29vb29vb29vb29vb29vZgo=",
   "ContentType" : "text/plain",
   "Hash" : "5la2x0TvvH7eaL1pcCN9QBVzbdULDdmld6HyGq3_SIg=",
   "FirstSeen" : 1411152910
}
```



`<base64_urlencode>` is the base64 representation of the sha256sum of a binary blob, file, or anything really.


##Testing

To see the service in action run (this requies httpie). 

```bash
$ http POST img.ahimsa.io/upload Content-Type:text/plain < create_table.sql

$ http GET img.ahimsa.io/eYQpbC_JWiz4tqZWaGmXGt8unK4A1Pgr6bhbdetWx4c=

$ http GET img.ahimsa.io/api/eYQpbC_JWiz4tqZWaGmXGt8unK4A1Pgr6bhbdetWx4c=
```

go test cases will come when more than 2 people use this!

