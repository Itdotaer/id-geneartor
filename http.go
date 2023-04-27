package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itdotaer/id-generator/generator"
	service2 "github.com/itdotaer/id-generator/service"
	"net"
	"net/http"
	"strconv"
	"time"
)

var GgeneratorService service2.GeneratorService

type allocResponse struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
	Id    int64  `json:"id"`
}

func handleAlloc(w http.ResponseWriter, r *http.Request) {
	var (
		resp     allocResponse = allocResponse{}
		err      error
		bytes    []byte
		business string
	)

	if err = r.ParseForm(); err != nil {
		goto RESP
	}

	if business = r.Form.Get("business"); business == "" {
		err = errors.New("need business param")
		goto RESP
	}

	for { // 跳过ID=0, 一般业务不支持ID=0
		if resp.Id, err = GgeneratorService.NextId(business); err != nil {
			goto RESP
		}
		if resp.Id != 0 {
			break
		}
	}

RESP:
	if err != nil {
		resp.Errno = -1
		resp.Msg = fmt.Sprintf("%v", err)
		w.WriteHeader(500)
	} else {
		resp.Msg = "success"
	}
	if bytes, err = json.Marshal(&resp); err == nil {
		w.Write(bytes)
	} else {
		w.WriteHeader(500)
	}
}

func StartServer() error {
	GgeneratorService = service2.NewGeneratorService()
	mux := http.NewServeMux()
	mux.HandleFunc("/alloc", handleAlloc)

	srv := &http.Server{
		ReadTimeout:  time.Duration(generator.GConf.HttpReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(generator.GConf.HttpWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(generator.GConf.HttpPort))
	if err != nil {
		return err
	}
	return srv.Serve(listener)
}
