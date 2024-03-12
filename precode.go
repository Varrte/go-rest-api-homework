package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...

// Обработчик для получения всех задач
func getTasks(res http.ResponseWriter, req *http.Request) {
	// сериализуем данные из слайса tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	// в заголовок записываем тип контента, у нас это данные в формате JSON
	res.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	res.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	res.Write(resp)
}

// Обработчик для отправки задачи на сервер
//{"id": "3","description":"111", "note":"222", "applications":["333","444","555"]}

func postTask(res http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	var addTask Task

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &addTask); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// addTask.ID это ключ, следовательно мы
	// добавляем к мапе tasks новый ключ и значение

	tasks[addTask.ID] = addTask
	res.Header().Set("Content-Type", "aplication/json")
	res.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID
func getTaskId(res http.ResponseWriter, req *http.Request) {
	//chi.URLParam()возвращает Url адрес из http.Request
	id := chi.URLParam(req, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(res, "Задача не найдена", http.StatusNoContent)
		return
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

// Обработчик удаления задачи по ID
func deleteTaskId(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	_, ok := tasks[id]
	if !ok {
		http.Error(res, "Задача не найдена", http.StatusNoContent)
		return
	}
	delete(tasks, id)

	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}
func main() {
	r := chi.NewRouter()
	// здесь регистрируйте ваши обработчики
	// ...
	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTaskId)
	r.Delete("/tasks/{id}", deleteTaskId)
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
