package user

import (
    _ "fmt"
    "net/http"
    "strconv"

    comp "wiesel/pb175/components"
    ut "wiesel/pb175/cmd/utility"

    _ "github.com/a-h/templ"
)


func Profile(st ut.GSP) ut.Response {
    return func(w http.ResponseWriter, r *http.Request) {
        user := ut.GetClientID(r)
        if user > -1 {
            http.SetCookie(w, ut.NewSession(user))
        }

        client := ut.GetUser(st, user)

        id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
            return
        }
        owner, err := st.DBH.GetUserById(id)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
            return
        }
        curr_page := comp.All
        if client.ID == owner.ID {
            curr_page = comp.ProfileP
        }
        comp.Page(
            comp.Profile(owner, client, ""),
            client,
            curr_page, 
        ).Render(r.Context(), w)
    }
}

func GetAllUsers(st ut.GSP) ut.Response {
    return func(w http.ResponseWriter, r *http.Request) {
        user := ut.GetClientID(r)
        if user > -1 {
            http.SetCookie(w, ut.NewSession(user))
        }

        client := ut.GetUser(st, user)

        if !client.IsAdmin {
            w.WriteHeader(http.StatusForbidden)
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
            return
        }
        users, _ := st.DBH.GetUsers()
        comp.Page(comp.Users(users, client), client, comp.All).Render(r.Context(), w)
    }
}

