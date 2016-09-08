package phpim

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func (im *IM) connectCallback(c *connection, r *http.Request) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if b, ok := e.(BadRequest); ok {
				err = b
			} else {
				panic(e)
			}
		}
	}()

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

	rs := CallbackResponse{}
	err = json.NewDecoder(resp.Body).Decode(&rs)
	if err != nil {
		return err
	}

	if rs.Ret != 0 {
		return errors.New("callback return fail.")
	}

	c.Id = rs.Id

	im.conns.Add(rs.Id, c)

	im.serveAction(rs.Actions)

	return nil
}

func (im *IM) msgCallback(c *connection, msg string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if b, ok := e.(BadRequest); ok {
				err = b
			} else {
				panic(e)
			}
		}
	}()

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

	im.serveAction(r.Actions)

	return nil
}

func (im *IM) disconnectCallback(c *connection) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if b, ok := e.(BadRequest); ok {
				err = b
			} else {
				panic(e)
			}
		}
	}()

	v := url.Values{}
	v.Set("act", "disconnect")
	v.Set("id", c.Id)

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

	for _, room := range c.rooms {
		room.Del(c)
	}

	im.conns.Del(c.Id)

	return nil
}
