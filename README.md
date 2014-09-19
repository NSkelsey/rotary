Rotary
======

A small go server that functions as a hash to file blob store. Storage is actually done in sqlite!


Routes
======

POST /upload
Upload some object via an http post. Returns 201 on success and greater than 400 on failure.

GET /<base64_urlencode> 
Returns the raw file with content-type set to conttype

GET /api/<base64_urlencode>
Returns a json object that looks like:
```


```




<base64_urlencode> is the base64 representation of the sha256sum of a binary blob, file, or anything really.
