package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"

	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	//	"github.com/h2non/imaginary"
)

func sendRequestToImaginary() {
	//	res, err := http.Post()
	resp, err := http.Get("imaginary:9000/form")
	if err != nil {
		//replace with a legit error later
		fmt.Println("trouble reaching imaginary server")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
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

	//filename format: uniquefilename_x_y.jpg
	resizedFilenameAndExt := strings.Split(resizedFilename,".")
	tokenizedName := strings.Split(resizedFilenameAndExt[0], "_")

	fmt.Println(tokenizedName)

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
		fileExt := fileInfo.Name();
		c.Response().Header().Set(echo.HeaderContentType, resizedFilename + "." + fileExt)
	  	return c.File("static/resized/" + resizedFilename)
	}

	return c.String(http.StatusOK, "temporary")
}

func main() {
	e := echo.New()

	e.GET("/static/resized/:resizedImageName", imgResizeHandler)
	e.Static("/static", "./static")

	e.Logger.Fatal(e.Start(":1323"))
}
