package server

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/nats-io/nats.go"
	"github.com/pion/ion-sfu/pkg/sfu"
	"github.com/pion/webrtc/v3"
	"log"
)

const (
	NatsChannel = "sfu."
)

// Join message sent when initializing a peer connection
type Join struct {
	SID    string                    `json:"sid"`
	UID    string                    `json:"uid"`
	Token  string                    `json:"token"`
	Offer  webrtc.SessionDescription `json:"offer"`
	Config sfu.JoinConfig            `json:"config"`
}

// Negotiation message sent when renegotiating the peer connection
type Negotiation struct {
	Uid  string                    `json:"uid"`
	Desc webrtc.SessionDescription `json:"desc"`
}

// Trickle message sent when renegotiating the peer connection
type Trickle struct {
	Uid       string                  `json:"uid"`
	Target    int                     `json:"target"`
	Candidate webrtc.ICECandidateInit `json:"candidate"`
}

func (n *NatsSignal) StartListen() {
	n.On("msg", func(msg *nats.Msg) {
		message := msg.Data
		room := getRoom(msg.Subject)
		replyid := getUid(msg.Subject)
		method := getMethod(msg.Subject)
		//token := getToken(msg.Subject)
		//t, err := n.ParseToken(token)
		//if err != nil {
		//	fmt.Println("jwt authorized! ", err)
		//	return
		//}
		//if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		//	fmt.Println("token: ", claims["sub"])
		//} else {
		//	fmt.Println("unauthorized")
		//	return
		//}
		fmt.Println("room: ", room, "method: ", method, "message: ", string(message))
		switch method {
		case "join":
			var join Join
			err := json.Unmarshal(message, &join)
			if err != nil {
				log.Fatal("error pasing join: ", err)
				break
			}

			token := join.Token
			t, err := n.ParseToken(token)
			if err != nil {
				log.Fatal("unauthorized", err)
				break
			}

			if _, ok := t.Claims.(jwt.MapClaims); !ok && !t.Valid {
				log.Fatal("unauthorized")
				break
			}

			peer := sfu.NewPeer(n.sfu)
			// when peer trigger onIceCandidate it sends back to biz
			peer.OnIceCandidate = func(candidate *webrtc.ICECandidateInit, target int) {
				message := Trickle{
					Target:    target,
					Candidate: *candidate,
				}
				resp, err := json.Marshal(message)
				if err != nil {
					log.Fatal("error marshal reply candidate: ", err)
					return
				}
				subject := NatsChannel + n.config.ServerName.String() + ".candidate"
				if err := n.natsconn.Publish(subject, resp); err != nil {
					log.Fatal("error publish candidate: ", err)
					return
				}
			}
			// when peer trigger OnOffer it will send back to biz
			peer.OnOffer = func(o *webrtc.SessionDescription) {
				offer, err := json.Marshal(o)
				if err != nil {
					log.Fatal("error marshal offer: ", err)
					return
				}
				subject := NatsChannel + n.config.ServerName.String() + ".offer"
				if err := n.natsconn.Publish(subject, offer); err != nil {
					log.Fatal("error publish offer: ", err)
				}
			}

			err = peer.Join(join.SID, join.UID)
			if err != nil {
				log.Fatal("failed to join: ", err)
				break
			}
			n.peers[join.UID] = peer

			answer, err := peer.Answer(join.Offer)
			if err != nil {
				log.Fatal("failed to generate answer: ", err)
			}

			message := Negotiation{
				Uid: peer.ID(),
				Desc: webrtc.SessionDescription{
					Type: answer.Type,
					SDP:  answer.SDP,
				},
			}

			resp, err := json.Marshal(message)
			if err != nil {
				log.Fatal("error marshal Answer: ", err)
				break
			}
			subject := NatsChannel + n.config.ServerName.String() + ".answer"
			if err := n.natsconn.Publish(subject, resp); err != nil {
				log.Fatal("error publish answer: ", err)
				break
			}

		case "offer":
			var negotiation Negotiation
			subject := NatsChannel + n.config.ServerName.String() + ".answer"
			err := json.Unmarshal(message, &negotiation)
			if err != nil {
				log.Fatal("connect: error parsing offer", err)
				break
			}

			p, found := n.peers[negotiation.Uid]
			if !found {
				//err
				log.Fatal("offer uid not found")
				break
			}
			answer, err := p.Answer(negotiation.Desc)
			if err != nil {
				log.Fatal("error unmarshal answer: ", err)
				break
			}
			message := Negotiation{
				Uid: p.ID(),
				Desc: webrtc.SessionDescription{
					Type: answer.Type,
					SDP:  answer.SDP,
				},
			}
			resp, err := json.Marshal(message)
			if err != nil {
				log.Fatal("error marshal answer: ", err)
			}

			if err := n.natsconn.Publish(subject, resp); err != nil {
				log.Fatal("failed to publish answer: ", err)
			}

		case "answer":
			var negotiation Negotiation
			err := json.Unmarshal(message, &negotiation)
			if err != nil {
				log.Fatal("fail to unmarshal answer message: ", err)
				break
			}
			p, found := n.peers[negotiation.Uid]
			if !found {
				log.Fatal("answer uid not found")
				break
			}
			err = p.SetRemoteDescription(negotiation.Desc)
			if err != nil {
				log.Fatal("failed to set remove discription: ", err)
				break
			}

		case "trickle":
			var trickle Trickle
			err := json.Unmarshal(message, &trickle)
			if err != nil {
				log.Fatal("failed to unmarshal trickle message: ", err)
				break
			}
			p, found := n.peers[trickle.Uid]
			if !found {
				log.Fatal("cannot find user: ")
				break
			}

			err = p.Trickle(trickle.Candidate, trickle.Target)
			if err != nil {
				log.Fatal("failed to trickle: ", err)
				break
			}

		case "leave":
			p, found := n.peers[replyid]
			if !found {
				log.Fatal("trickle uid not found")
				break
			}
			p.Close()
		}
	})
}
