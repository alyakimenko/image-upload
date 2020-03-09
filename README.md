### How to run
To run the app just execute `make` command, to build and create executable - `make build`, and to run tests - `make test`.
Server will start on `localhost:8080` after that.
### Description
The server has only one route. `POST /upload`.
It can take image (or multiple images) as:
- multipart/form-data
- base64 encoded string represents in JSON
- list of urls

After upload each image will be resized to `100x100`, and stored in `./downloaded/` directory (by default).
### Examples
**Multipart/form-data:**

POST request on `http://localhost:8080/upload` with `multipart/form-data` Content-Type.
Example with curl:
```shell script
curl -F files=@1.jpg -F files=@2.jpg http://localhost:8080/upload
```

**BASE64 JSON:**

POST request on `http://localhost:8080/upload` with `application/json` Content-Type.
Example of allowed JSON:
```json
{
  "data": [
    "/9j/4AAQSkZJRgABAQAAAQABAAD/4QDeRXhpZgA...",
    "zk9ODI8LjM0Mv/bAEMBCQkJDAsMGA0NGDIhHCEyMjIy..."
  ]
}
```
Each BASE64 encoded string will be decoded to image, if possible.

**URLs:**

POST request on `http://localhost:8080/upload` with `text/plain` Content-Type.
Each url must be splitted up by newline symbol (`\n`).
For instance:
```text
https://i.picsum.photos/id/237/200/300.jpg
https://i.picsum.photos/id/245/200/300.jpg
https://i.picsum.photos/id/32/200/300.jpg
```
