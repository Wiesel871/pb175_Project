package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
    "bufio"

	cmd "wiesel/pb175/cmd"
	_ "wiesel/pb175/components"
	data "wiesel/pb175/database"
	gst "wiesel/pb175/state"

	_ "github.com/a-h/templ"
)

func main() {
    dbh, err := data.InitDB("bazos.db")
    if err != nil {
        return
    }


    st := gst.GlobalState {
        DBH: dbh,
        Anonym: gst.GetAnonym(),
    }

    mux := http.NewServeMux()
    cmd.SetupUserHandler(mux, &st)
    st.SRV = &http.Server {
        Addr: ":8090",
        Handler: mux,
        ErrorLog: log.Default(),
    }


    file, err := os.Open("admin.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewReader(file)
    name, _, _ := scanner.ReadLine()
    email, _, _ := scanner.ReadLine()
    password, _, _ := scanner.ReadLine()

    id := int64(0)
    admin, _ := data.NewUser(id, string(name), string(email), string(password))
    st.DBH.InsertUser(admin)
    st.DBH.Promote(id)
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan

        shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
        defer shutdownRelease()

        if err := st.SRV.Shutdown(shutdownCtx); err != nil {
            log.Fatalf("HTTP shutdown error: %v", err)
        }
    }()

    if err := st.SRV.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        fmt.Printf("err: %v\n", err)
    }
    fmt.Println("Shutdown")
    st.DBH.DB.Close()
}
