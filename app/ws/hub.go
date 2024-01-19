package ws

type ConsultationSession struct {
	Id        int64             `json:"id"`
	DoctorId  int64             `json:"doctor_id"`
	PatientId int64             `json:"patient_id"`
	Clients   map[int64]*Client `json:"clients"`
}

type Hub struct {
	ConsultationSessions map[int64]*ConsultationSession
	Register             chan *Client
	Unregister           chan *Client
	Broadcast            chan *Message
}

func NewHub() *Hub {
	return &Hub{
		ConsultationSessions: make(map[int64]*ConsultationSession),
		Register:             make(chan *Client),
		Unregister:           make(chan *Client),
		Broadcast:            make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, isRoomExist := h.ConsultationSessions[client.RoomId]; isRoomExist {
				r := h.ConsultationSessions[client.RoomId]

				if _, isClientExist := r.Clients[client.Id]; !isClientExist {
					r.Clients[client.Id] = client
				}
			}
		case client := <-h.Unregister:
			if _, isRoomExist := h.ConsultationSessions[client.RoomId]; isRoomExist {
				if _, isClientExist := h.ConsultationSessions[client.RoomId].Clients[client.Id]; isClientExist {
					if len(h.ConsultationSessions[client.RoomId].Clients) != 0 {
						h.Broadcast <- &Message{
							Content: ConsultationMessage{
								Message: "A user has left the room chat",
							},
							RoomId: client.RoomId,
						}
					}

					delete(h.ConsultationSessions[client.RoomId].Clients, client.Id)
					close(client.Message)
				}
			}
		case message := <-h.Broadcast:
			if _, isRoomExist := h.ConsultationSessions[message.RoomId]; isRoomExist {

				for _, client := range h.ConsultationSessions[message.RoomId].Clients {
					client.Message <- message
				}
			}
		}
	}
}
