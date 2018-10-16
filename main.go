package main

import (
	"fmt"

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

func imgDefaultHandler(c echo.Context) error {
	//e.GET("/:rawImagePath/:resizedImageName", imgDefaultHandler)

	//multiple images that have the same filename
	//to deal with multiple images that might have the same name, I tihnk it'd be worthwhile to hash the function. I know I'd have to worry about namespace collisions in a real-world application but

	//images of different aspect ratios and original sizes

	//run req validator function here and respond with error early if needed
	err := validateRequestContext(c)
	if err != nil {
		fmt.Println(err)
	}
	return nil
	//	c.JSON(http.StatusBadRequest, jsonErrorResPayload)
}

func main() {
	e := echo.New()

	e.GET("/:rawImagePath/:resizedImageName", imgDefaultHandler)
	e.GET("/:resizedImagePath/:resizedImageName", imgResizeHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
