//Пишем HTTP-сервис, который принимает входящие соединения с JSON-данными и обрабатывает их

package main

import (
	"encoding/json"
	"fmt"
	chi "github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
	"strconv"
  "structs"
  "helpers"
)

var S = Storage{users: make(map[int]*User), lastUserId: 0}

/* 1. Обработчик создания пользователя. У пользователя должны быть следующие поля: имя, возраст и массив друзей. Пользователя необходимо сохранять в мапу.
Данный запрос должен возвращать ID пользователя и статус 201*/

func Create(w http.ResponseWriter, r *http.Request) {

	var u User

	//берем и считываем body запроса
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//при ошибке напишем ошибку сервера 500 Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	//не забываем закрывать
	defer r.Body.Close()

	//десериализация JSON, передаем в функцию Unmarshal контент и указатель на пользователя
	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	//присваиваем id и складываем указатели на  user в storage
	UserId := AddNewUser(&u)

	//возвращаем ID пользователя и статус 201
	w.WriteHeader(http.StatusCreated)
	w.Write(([]byte(fmt.Sprintf("ID: %v  name: %v  age: %d  friend: %v\n", UserId, u.Name, u.Age, u.Friends))))
	w.Header().Set("Content-Type", "application/json")

	//записываем данные всех пользователей в виде JSON в файл
	InFile()
}

/*2. Обработчик, который делает друзей из двух пользователей. Например, если мы создали двух пользователей и нам вернулись их ID,
то в запросе мы можем указать ID пользователя, который инициировал запрос на дружбу, и ID пользователя, который примет инициатора в друзья.
Данный запрос должен возвращать статус 200 и сообщение «username_1 и username_2 теперь друзья».*/

func MakeFriends(w http.ResponseWriter, r *http.Request) {

	var Friends MakeFriendsRequest

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(content, &Friends); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	sourceId := Friends.SourceId
	targetId := Friends.TargetId

	u1, ok1 := S.IdExistenceCheck(sourceId)
	u2, ok2 := S.IdExistenceCheck(targetId)

	if !(ok1 && ok2) {
		w.WriteHeader(http.StatusBadRequest)
		r := MakeFriendsResponse{Message: "Wrong ID! There is no such ID"}
		response, _ := json.Marshal(r)
		w.Write(response)
		return
	}

	u1.Friends = append(u1.Friends, targetId)
	u2.Friends = append(u2.Friends, sourceId)

	for id, u := range S.users {
		fmt.Printf("ID: %v  name: %v  age: %d  friend: %v\n", id, u.Name, u.Age, u.Friends)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(([]byte(fmt.Sprintf("%v and %v are now friends", u1.Name, u2.Name))))
	w.WriteHeader(http.StatusOK)

	//записываем данные всех пользователей в виде JSON в файл
	InFile()
}

/*3. Обработчик, который принимает ID пользователя и удаляет его из хранилища, а также стирает его из массива friends у всех его друзей.
Данный запрос должен возвращать 200 и имя удалённого пользователя.*/

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	var UserForDelete DeleteUserRequest

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(content, &UserForDelete); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	IdForDelete := UserForDelete.TargetId
	deletedUser := S.users[IdForDelete]

	// удаляем ключ из map
	delete(S.users, IdForDelete)

	// удаляем значения из массива
	for _, j := range S.users {
		for m, k := range j.Friends {
			if k == IdForDelete {
				if len(j.Friends) > 1 {
					j.Friends = append(j.Friends[:m], j.Friends[m+1:]...)
				} else {
					j.Friends = nil
				}
			}
		}
	}

	for id, u := range S.users {
		fmt.Printf("ID: %v  name: %v  age: %d  friend: %v\n", id, u.Name, u.Age, u.Friends)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(([]byte(fmt.Sprintf("%v was deleted", deletedUser.Name))))
	w.WriteHeader(http.StatusOK)

	//записываем данные всех пользователей в виде JSON в файл
	InFile()
}

/*4. Обработчик, который возвращает всех друзей пользователя.
После /friends/ указывается id пользователя, друзей которого мы хотим увидеть.*/

func GetFriends(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	u := S.users[id]
	b, err := json.Marshal(u.Friends)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	fmt.Printf("Список друзей для %v\n", u)

	//записываем данные всех пользователей в виде JSON в файл
	InFile()
}

/*5. Сделайте обработчик, который обновляет возраст пользователя.
Запрос должен возвращать 200 и сообщение «возраст пользователя успешно обновлён».*/

func SetAge(w http.ResponseWriter, r *http.Request) {

	var uAge UserAge

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(content, &uAge); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	S.users[id].Age = uAge.Age

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Возраст пользователя обновлен"))

	//записываем данные всех пользователей в виде JSON в файл
	InFile()
}
