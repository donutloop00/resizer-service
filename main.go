package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
	"net/http"
)

func sendRequestToResizer(fileHandle *os.File, x int, y int) (string, error) {
	splitFileName := strings.Split(fileHandle.Name(), ".")
	tokenizedPathNoExt := strings.Split(splitFileName[0], "/")
	baseFileName := tokenizedPathNoExt[len(tokenizedPathNoExt)-1]
	imageType := splitFileName[len(splitFileName) - 1]
	contentType := "image/" + imageType

	widthAsStr := strconv.Itoa(x)
	heightAsStr := strconv.Itoa(y)

	resizeQueryStr := "width=" + widthAsStr + "&height=" + heightAsStr

	resp, err := http.Post("http://imaginary:9000/resize?"+resizeQueryStr, contentType, fileHandle)
	if err != nil {
		errorMsg := "Error communicating with image resize server."
		return "", echo.NewHTTPError(http.StatusInternalServerError, errorMsg)
	}
	defer resp.Body.Close()
	fmt.Println(resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorMsg := "Error reading image response from resizer service."
		return "", echo.NewHTTPError(http.StatusInternalServerError, errorMsg)
	}

	fmt.Println(fileHandle.Name())

	pathOfResizedImage := "static/resized/" + baseFileName + "_" + widthAsStr + "_" + heightAsStr + "." + imageType
	err = ioutil.WriteFile(pathOfResizedImage, body, 0666)
	if err != nil {
		errorMsg := "Error writing resized image to server disk."
		return "", echo.NewHTTPError(http.StatusInternalServerError, errorMsg)
	}

	return pathOfResizedImage, nil
}

func checkValidDimensions(x, y int) bool {
	if x > 0 && y > 0 {
		return true
	} else {
		return false
	}
}

func imgResizeHandler(c echo.Context) error {
	//check if the image has already been resized before
	resizedFilename := c.Param("resizedImageName")

	//filename format: uniquefilename-x-y.jpg
	resizedFilenameAndExt := strings.Split(resizedFilename, ".")
	fileExt := resizedFilenameAndExt[len(resizedFilenameAndExt)-1]
	tokenizedName := strings.Split(resizedFilenameAndExt[0], "-")

	reqWidth, err := strconv.Atoi(tokenizedName[1])
	if err != nil {
		errorMsg := "Error converting width to an int"
		return echo.NewHTTPError(http.StatusBadRequest, errorMsg)
	}

	reqHeight, err := strconv.Atoi(tokenizedName[2])
	if err != nil {
		errorMsg := "Error converting height to an int"
		return echo.NewHTTPError(http.StatusBadRequest, errorMsg)
	}

	dimensionsAreValid := checkValidDimensions(reqWidth, reqHeight)

	if !dimensionsAreValid {
		errorMsg := "One or both height and width dimensions are invalid"
		return echo.NewHTTPError(http.StatusBadRequest, errorMsg)
	}

	// see if the image has been resized before
	if fileInfo, err := os.Stat("static/resized/" + resizedFilename); !os.IsNotExist(err) {
		//trusting the server to get the file extension right instead of the user's requested filename
		fileExt := fileInfo.Name()
		c.Response().Header().Set(echo.HeaderContentType, resizedFilename+"."+fileExt)
		return c.File("static/resized/" + resizedFilename)
	}

	// check if the source image exists. if not, send a 404 error
	pathOfSourceImage := "static/source/" + tokenizedName[0] + "." + fileExt
	if _, err := os.Stat(pathOfSourceImage); err != nil {
		if os.IsNotExist(err) {
			errorMsg := "Source image does not exist for resizing."
			return echo.NewHTTPError(http.StatusNotFound, errorMsg)
		}
	} else {

		fileAsIOReader, err := os.Open(pathOfSourceImage)
		if err != nil {
			errorMsg := "Error opening image on server."
			return echo.NewHTTPError(http.StatusNotFound, errorMsg)
		}

		pathOfResizedImage, err := sendRequestToResizer(fileAsIOReader, reqWidth, reqHeight)
		if err != nil {
			//all errors are already formatted as echo.Context errors
			return err
		}
		return c.File(pathOfResizedImage)
	}

	return nil
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/static/resized/:resizedImageName", imgResizeHandler)
	e.Static("/static", "./static")

	e.Logger.Fatal(e.Start(":1323"))
}
