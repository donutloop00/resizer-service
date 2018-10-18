
# resizer-service 

An image resizer service in Go.

I had originally intended to use the `scratch` docker image to keep the Docker builds lean, but ran into problems trying to get other Linux shell commands working in them, so I switched over to the Golang 1.11 distro which depends on Debian Stretch.

To deal with unique filenames, a separate service can implement an uploader that converts the uploaded filename into a hash string using whatever hashing algorithm of choice you'd like. That way, whatever files are in the `static/source` directory are already unique both in the form of filename and content.

To run:

Move a test image to in the filename to the `static/source` directory. Make sure there are no '-' delimiters in the source image name used.
Also make sure that the Golang `dep` tool is installed so that `dep ensure` will work correctly.

```
mkdir -p static/source
mkdir -p static/resized
dep ensure
make
docker-compose build
docker-compose up
```

##Routes:

### `/static/source/baseImageFilename`

This route is simply implemented as a static folder. 

Responses:

+ 200 OK with header { Content-Type: image/jpg, image/png, image/gif.}
+ 404 if File not Found 

### `/static/resized/baseImageFilename-width-height.jpg`

baseImageFilename must not have any `-` or `.` characters present within it because those are being used as delimiters to determine filename, width, height, and file extension.

+ 200 OK with header { Content-Type: image/jpg, image/png, image/gif.}
+ 400 Bad Request if the dimensions don't make sense (0 or negative dimensions along an axis)
+ 404 if source image does not exist 
+ 500 if HTTP requests to h2non/imaginary container fail or if the resized image cannot be saved to the `static/resized` folder 

## Nice-to-Have's If This Were An Actual Service:

I think long-term, it'd be worth checking the logs of a service like this to see what are the most commonly-requested sizes. That and adding a more convenient logger format if the Echo middleware version is incovenient.

## Shortcuts

So some things I had cut out in the interest of getting this done slightly sooner:

+ Setting up proper environment variable flags for the original app Docker container. Normally I'd have an `.env` file and `godotenv` to store any API keys or secrets but I didn't see a point in setting up the secrets for `h2non/imaginary` when the service is this small
+ There is a bug in the most recent Docker for MacOS where DNS names weren't resolving correctly, so it did take me some more time to figure that out (since the app service sends an HTTP request to `imaginary:9000` but this request fails in Docker Compose when local DNS isn't resolving correctly. I fixed this by updating to the latest Docker version.
+ Error handling for additional edge cases, such as when a user enters a height/width pair that does not maintain aspect ratio. Even though the documentation for `imaginary` says it maintains aspect ratio, it is possible to submit a request that skews/shears the resized image (e.g. requesting 10x1000 on an image that is originally 4:6) if we send a request with only width defined for resizing.
+ Additional unit tests.
+ Covering additional edge cases with the filename. 
+ There is also a bug to fix where if a user requests a resized image and receives it, and requests it again, there's no check on the server-side to verify that the dimensions of the file in the `static/resized` folder entirely match what was requested.

Even though I already had the project in working order with the Go API service running on `localhost` and running `h2non/imaginary` in a single Docker container, I decided to take the extra time to make sure I could get container-to-container communication working correctly with Docker Compose.
