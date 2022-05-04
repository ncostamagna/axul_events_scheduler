package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ncostamagna/alertzy-sdk/alertzy"
	c "github.com/ncostamagna/streetflow/client"
)

const (
	BASE          = ""
	USERNAME      = ""
	PASSWORD      = ""
	ALERTZY_TOKEN = ""
)

func main() {

	var err error

	id, token, err := userReq()
	if err != nil {
		panic(err)
	}
	fmt.Println(id, token)

	contacts, err := contReq(id, token)
	if err != nil {
		panic(err)
	}
	fmt.Println(contacts)

	cTrans := alertzy.NewClient("https://alertzy.app", ALERTZY_TOKEN)

	for _, c := range contacts {
		err = cTrans.Send(fmt.Sprintf("Hoy es el cumple de %s", c), fmt.Sprintf("Hoy es el cunple de %s acordate de saludarlo en su dia", c), alertzy.Critical, "birthday", "", "", nil)
		fmt.Print(err)
	}

	os.Exit(0)

}

func userReq() (string, string, error) {

	client := c.RequestBuilder{
		Headers:        http.Header{},
		BaseURL:        BASE,
		ConnectTimeout: 5000 * time.Millisecond,
		LogTime:        true,
	}

	req := struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}{USERNAME, PASSWORD}

	reps := client.Post("/users/login", req)

	if reps.Err != nil {
		return "", "", reps.Err
	}

	if reps.StatusCode > 299 {
		return "", "", fmt.Errorf("code: %d, message: %s", reps.StatusCode, reps)
	}

	res := struct {
		Data struct {
			User struct {
				ID string `json:"id"`
			} `json:"user"`
			Token string `json:"token"`
		} `json:"data"`
	}{}

	if err := json.Unmarshal(reps.Bytes(), &res); err != nil {
		return "", "", err
	}

	return res.Data.User.ID, res.Data.Token, nil

}

func contReq(id, token string) ([]string, error) {

	header := http.Header{}
	header.Set("Authorization", token)

	client := c.RequestBuilder{
		Headers:        header,
		BaseURL:        BASE,
		ConnectTimeout: 5000 * time.Millisecond,
		LogTime:        true,
	}

	reps := client.Get(fmt.Sprintf("/contacts?birthday=0&userid=%s", id))

	if reps.Err != nil {
		return nil, reps.Err
	}

	if reps.StatusCode > 299 {
		return nil, fmt.Errorf("code: %d, message: %s", reps.StatusCode, reps)
	}

	res := struct {
		Data []struct {
			FirstName string `json:"firstname"`
			LastName  string `json:"lastname"`
		} `json:"data"`
	}{}

	if err := json.Unmarshal(reps.Bytes(), &res); err != nil {
		return nil, err
	}

	var contacts []string

	for _, c := range res.Data {
		contacts = append(contacts, fmt.Sprintf("%s %s", c.FirstName, c.LastName))
	}

	return contacts, nil

}
