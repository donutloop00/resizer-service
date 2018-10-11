package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	//"github.com/labstack/echo"
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

func imgDefaultHandler(c echo.Context) {
	//run req validator function here and respond with error early if needed

}

func imgResizeHandler(c echo.Context) {
	//run req validator function here and respond with error early if needed

}

func main() {
	e := echo.New()

	e.GET("/path1", imgDefaultHandler)
	e.GET("/path2", imgResizeHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
