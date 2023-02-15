package helpers

// функция записи указателя в массив и изменения ID
func AddNewUser(u *User) int {
	S.users[S.lastUserId] = u
	S.lastUserId++
	return S.lastUserId - 1
}

// функция проверки значений на существование
func (s *Storage) IdExistenceCheck(id int) (*User, bool) {
	u, ok := s.users[id]
	return u, ok
}

// функция записи данных JSON в файл
func InFile() {
	file, err := os.OpenFile("user.json", os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for id, u := range S.users {
		response, err := json.Marshal(CreateUserResponse{Id: id, Name: u.Name, Age: u.Age, Friends: u.Friends})
		if err != nil {
			fmt.Println(err)
			return
		}
		file.Write(response)
	}
}
