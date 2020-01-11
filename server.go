package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "mkserver/pb"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12345"
	dbname   = "suser"
)

func RunServer() {

	r := mux.NewRouter()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	fmt.Println(psqlInfo)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("errrrror", err)
		return
	}
	defer db.Close()
	fmt.Println("connected")
	db.CreateTable(&pb.MyMsgORM{})
	r.HandleFunc("/insertmsg", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Println("parsed form")
		user, err := pb.DefaultCreateMyMsg(context.Background(), &pb.MyMsg{Msg1: r.FormValue("msg1"), Msg2: r.FormValue("msg2"), Msg3: r.FormValue("msg3")}, db)
		if err != nil {
			panic(err)
		}
		fmt.Println("msg created with info", user)

	})
	r.HandleFunc("/displaymsg", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Println("parsed form")
		var msg pb.MyMsgORM
		err := db.Where("msg1=?", r.FormValue("searchmsg")).First(&msg).Error

		user, err := pb.DefaultReadMyMsg(context.Background(), &pb.MyMsg{Id: msg.Id}, db)
		if err != nil {
			fmt.Println("errrrrror", err)
			return
		}
		fmt.Println("msg read with info", user)

	})

	r.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("serverRunHIt success")
		return
	})
	http.ListenAndServe(":80", r)

}
