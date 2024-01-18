package ws

type Room struct {
	Id        int64             `json:"id"`
	DoctorId  int64             `json:"doctor_id"`
	PatientId int64             `json:"patient_id"`
	Clients   map[int64]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[int64]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[int64]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, isRoomExist := h.Rooms[client.RoomId]; isRoomExist {
				r := h.Rooms[client.RoomId]

				if _, isClientExist := r.Clients[client.Id]; !isClientExist {
					r.Clients[client.Id] = client
				}
			}
		case client := <-h.Unregister:
			if _, isRoomExist := h.Rooms[client.RoomId]; isRoomExist {
				if _, isClientExist := h.Rooms[client.RoomId].Clients[client.Id]; isClientExist {
					if len(h.Rooms[client.RoomId].Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  "user left the chat",
							RoomId:   client.RoomId,
						}
					}

					delete(h.Rooms[client.RoomId].Clients, client.Id)
					close(client.Message)
				}
			}
		case message := <-h.Broadcast:
			if _, isRoomExist := h.Rooms[message.RoomId]; isRoomExist {

				for _, client := range h.Rooms[message.RoomId].Clients {
					client.Message <- message
				}
			}
		}
	}
}
