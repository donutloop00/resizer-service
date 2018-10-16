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

func checkValidDimensions(x, y int) bool {
	if x > 0 && y > 0 {
		return true
	} else {
		return false
	}
}

func main() {
	e := echo.New()

	e.Static("/static", "./static")

	e.Logger.Fatal(e.Start(":1323"))
}
