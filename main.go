/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   main.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: asyed <asyed@student.42.fr>                +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2018/05/22 16:54:47 by asyed             #+#    #+#             */
/*   Updated: 2018/05/24 14:02:12 by asyed            ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"fmt"
	"net"
	"bytes"
	"strings"
	"net/url"
	"errors"
	"C"
	"encoding/hex"
	// "encoding/binary"
	// "unsafe"
	// "github.com/anacrolix/torrent"
	// "reflect"
)

type Magnet struct {
	values	*url.Values; // values	*map[string][]string;
	// dn	string;
	// xl	uint64;
	// xt	string;
	// as	string;
	// xs	string;
	// kt	string;
	// mt	string;
	// tr	[]string;
};

type ipv4_announce_request struct {
	connection_id	uint64;
	action			uint32;
	transaction_id	uint32;
	info_hash		[20]byte;
	peer_id			[20]byte;
	downloaded		uint64;
	left			uint64;
	uploaded		uint64;
	event			uint32;
	ip_addr			uint32;
	key				uint32;
	num_want		int32;
	port			uint16;
	pad				[2]byte;
};

type Tracker struct {
	announce_request	ipv4_announce_request;
	// connection_id	uint64;
	// info_hash		[20]byte;

}

func (magnet *Magnet) parse_magnet(query *url.Values) (error) {
	if query == nil {
		return (errors.New("Query is empty"));
	}
	magnet.values = query;
	return (nil);
}

/*
** I should use a worker pool for the downloading and then return
** a link/stream to the file. 
**
** Do we want the traffic for the stream to tunnel through this service
** as well or do we want to have it mandatorily be on the torrent/downloading
** server? (This server.)
**
**
** Make this so that we can properly escape the hash.
*/

func (magnet *Magnet) download_torrent() (error) {
	var addr			*net.UDPAddr;
	var parse			*url.URL;
	var conn			*net.UDPConn;
	var tracker			Tracker;
	var err				error;

	if parse, err = url.ParseRequestURI((*magnet.values)["tr"][0]); err != nil {
		return (err);
	}
	if addr, err = net.ResolveUDPAddr("udp", parse.Host); err != nil {
		return (err);
	}
	if conn, err = net.DialUDP("udp", nil, addr); err != nil {
		return (err);
	}
	defer conn.Close();
	if ((*magnet.values)["xt"][0] == "") {
		return errors.New("Empty hash!");
	}
	(*magnet.values)["xt"][0] = strings.TrimPrefix((*magnet.values)["xt"][0], "urn:");
	(*magnet.values)["xt"][0] = strings.TrimPrefix((*magnet.values)["xt"][0], "btih:");
	fmt.Println("Size of hash = ", len((*magnet.values)["xt"][0]));
	if ret, err := hex.Decode(tracker.announce_request.info_hash[:20], []byte((*magnet.values)["xt"][0])); err != nil || ret != 20 {
		if (err == nil) {
			fmt.Println(ret, "instead of 20");
			return errors.New("Didn't copy the correct size");
		}
		return err;
	}
	if err = tracker.handshake(conn);err != nil {
		return (err);
	}
	fmt.Println(conn.LocalAddr().String());
	fmt.Printf("Hash = \"%s\"\n", string(tracker.announce_request.info_hash[:20]))
	fmt.Printf("Connection ID = %d\n", tracker.announce_request.connection_id);
	// for ;; {
	// 	var	ret int;
	// 	if ret, _, err = conn.ReadFromUDP(buffer); err != nil || ret != 16 {
	// 		return (err);
	// 	}
	// 	if (ret != 16) {
	// 		fmt.Println("Didn't read 16 ", ret);
	// 	}
	// 	fmt.Println("Read" + string(buffer[:ret]));
	// }
	return (nil);
}

/*
** The port number that the client is listening on.
** Ports reserved for BitTorrent are typically 6881-6889.
** Clients may choose to give up if it cannot establish a port within this range.
**
** Create a struct to make functions specific to this listener - IE: storing the status
** of the file thus handleConnection is specific to file downloads
*/

func handleConnection(conn net.Conn) {
	var buf	bytes.Buffer;

	buf.Grow(1024);
	if _, err := conn.Read(buf.Bytes());err != nil {
		fmt.Println("Error:", err.Error());
	}
	conn.Close();
	fmt.Println(buf.Bytes());
}

func torrentListen() (error) {
	ln, err := net.Listen("tcp", ":6881");
	if err != nil {
		return (err);
	}
	fmt.Println("Listening on", ln.Addr());
	for {
		conn, err := ln.Accept();
		if err != nil {
			return (err);
		}
		go handleConnection(conn);
	}
	return (nil);
}

func main() {
	const sample string = "magnet:?xt=urn:btih:3f2a441e2c4b84f25e44403328aeffc432a15ae7&dn=Ready+Player+One.2018.HDRip.X264.AC3-EVO%5BN1C%5D&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969";
	var parse	*url.URL;
	var query	url.Values;
	var err		error;
	var magnet	Magnet;

	go func() {
		if err = torrentListen();err != nil {
			panic(err);
		}
	}()
	if parse, err = url.ParseRequestURI(sample); err != nil {
		fmt.Printf("Error: %s\n", err);
		return ;
	}
	query = parse.Query();
	if err = magnet.parse_magnet(&query); err != nil {
		fmt.Printf("Error: \"%s\"\n", err);
		return ;
	}
	fmt.Printf("%p\n", query);
	for i := 0; i < len((*magnet.values)["tr"]); i++ {
		fmt.Printf("tr[%d] = \"%s\"\n", i, (*magnet.values)["tr"][i]);
	}
	if err = magnet.download_torrent(); err != nil {
		fmt.Printf("Error: %s\n", err);
		return ;
	}
	fmt.Printf("OK!\n");
	// fmt.Printf("%s", sample);
}