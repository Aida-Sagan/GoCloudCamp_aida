package main

import (
	"fmt"
	"time"
)

type Song struct {
	nameSong string
	duration time.Duration
}

type Playlist struct {
	head        *Node
	tail        *Node
	current     *Node
	playChannel chan bool
}

type Node struct {
	song *Song
	prev *Node
	next *Node
}

func MyPlaylist() *Playlist {
	return &Playlist{
		head:        nil,
		tail:        nil,
		current:     nil,
		playChannel: make(chan bool),
	}
}

func (p *Playlist) Play() {
	p.playChannel <- true
}

func (p *Playlist) AddSong(name *Song) {
	newSong := &Node{song: name, prev: p.tail}

	if p.head != nil {
		p.tail.next = newSong
		newSong.prev = p.tail
		p.tail = newSong
	} else {
		p.head = newSong
		p.tail = newSong
	}

	fmt.Printf("%s is added\n", name)
	fmt.Println()
}

func (p *Playlist) Pause() {
	p.playChannel <- false

}

func (p *Playlist) Next() {
	//воспроизводим первую песню, если список пуст
	if p.current == nil {
		p.current = p.head

	} else if p.current.next != nil {
		p.current = p.current.next
		fmt.Printf("Song: %s switched", p.current.song.nameSong)

	}
	p.current = p.head
	fmt.Println()
	p.playSong()

}

func (p *Playlist) Prev() {
	if p.current == nil {
		p.current = p.tail
	} else if p.current.prev != nil {
		p.current = p.current.prev
	}
	p.current = p.tail

	p.playSong()
}

func (p *Playlist) playSong() {
	if p.current == nil {
		return
	}
	//горутина для воспроизведения песни и отправляет смс через канал playChannel, когда песня завершится
	go func() {
		p.playChannel <- true
		time.Sleep(p.current.song.duration)
		p.Next()
	}()
}

func (p *Playlist) Start() {
	for {
		play := <-p.playChannel
		if p.current == nil {
			continue
		}
		if play {
			fmt.Printf("is playing: %s\n", p.current.song.nameSong)
			time.Sleep(p.current.song.duration)
			p.Next()
		} else {
			fmt.Printf("is Paused: %s\n", p.current.song.nameSong)
		}
	}
}

func main() {
	pl := MyPlaylist()

	// Добавление песен в плейлист
	song1 := &Song{nameSong: "happy song 1", duration: 2 * time.Second}
	song2 := &Song{nameSong: "sad song 2", duration: 3 * time.Second}
	//song3 := &Song{nameSong: "rock song 3", duration: 4 * time.Second}
	//song4 := &Song{nameSong: "pop song 4", duration: 5 * time.Second}

	pl.AddSong(song1)
	pl.AddSong(song2)
	//pl.AddSong(song3)
	//pl.AddSong(song4)

	//start our playlist
	go pl.Start()

	//play our playlist
	pl.playSong()

	//waiting for several sec
	time.Sleep(10 * time.Second)
	pl.playSong()

	//pause
	pl.Pause()

	//Next method
	pl.Next()

	pl.Pause()

	//Prev method
	pl.Prev()

	//waiting for ... sec
	time.Sleep(10 * time.Second)

	//close
	close(pl.playChannel)
	fmt.Println("Good bye, my friend!")
}
