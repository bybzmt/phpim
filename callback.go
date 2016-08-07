package phpim

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)


func (im *IM) connectCallback(c *conn) error {
	v := url.Values{}
	for key, ma := range r.URL.Query() {
		v.Add(key, ma[0])
	}

	v.Set("act", "connect")

	resp, err := http.PostForm(im.CallbackUrl, v)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	r := CallbackResponse{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return err
	}

	if r.Ret != 0 {
		return errors.New("callback return fail.")
	}

	c.Id = r.Id

	im.Conns.Add(r.Id, c)

	im.serveAction(ca.actions)

	return nil
}

func (im *IM) msgCallback(c *conn, msg string) {
	v := url.Values{}
	v.Set("act", "msg")
	v.Set("id", c.Id)
	v.Set("msg", msg)

	resp, err := http.PostForm(im.CallbackUrl, v)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	r := CallbackResponse{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return err
	}

	if r.Ret != 0 {
		return errors.New("callback return fail.")
	}

	im.serveAction(ca.actions)
}

func (im *IM) disconnectCallback(c *conn) {
	v := url.Values{}
	v.Set("act", "disconnect")
	v.Set("id", c.Id)

	resp, err := http.PostForm(im.CallbackUrl, v)
	if err != nil {
		log.Println("callback url error: ", err)
		return
	}

	defer resp.Body.Close()

	r := CallbackResponse{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return err
	}

	if r.Ret != 0 {
		log.Println(errors.New("callback return fail."))
	}

	for _, room := range c.Rooms {
		room.Del(c)
	}

	im.Conns.Del(c.Id, c)
}
