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

func main() {

	//	http.ListenAndServe(":8080")
}
