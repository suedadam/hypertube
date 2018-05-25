/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   announce.go                                        :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: asyed <asyed@student.42.fr>                +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2018/05/23 12:53:54 by asyed             #+#    #+#             */
/*   Updated: 2018/05/25 13:10:19 by asyed            ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"fmt"
	"errors"
	"net"
	"bytes"
	// "math"
	"encoding/binary"
)

func fuckit(size int, buffer bytes.Buffer) {
	for i := 0; i < size; i++ {
		fmt.Printf("\\x%02x", buffer.Bytes()[i]);
	}
	fmt.Printf("\n");
}

// const packetsize = 1024;

func (tracker *Tracker) ipv4_read_announce(conn *net.UDPConn) (error) {
	var response 	ipv4_announce_response;
	var buffer		*bytes.Buffer;
	var ret			int;
	var err			error;

	if buffer = bytes.NewBuffer(make([]byte, 2048)); buffer == nil {
		return (errors.New("Failed to initialzie buffer"));
	}
	if ret, err = conn.Read((*buffer).Bytes());err != nil || ret < 20 {
		if (err == nil) {
			fmt.Println("Read %d bytes (< 20)", ret);
			return (errors.New("Didn't read the whole packet"));
		}
		return (err);
	} else {
		fmt.Printf("Total return size = %d\n", ret);
	}
	response.action = binary.BigEndian.Uint32((*buffer).Bytes());
	if response.action != 1 {
		fmt.Println("Invalid action!? %d", response.action);
	} else {
		fmt.Println("Action okay!");
	}
	if response.transaction_id = binary.BigEndian.Uint32((*buffer).Bytes()[4:]); response.transaction_id != tracker.announce_request.transaction_id {
		fmt.Println("Transaction ID is incorrect %d instead of %d", response.transaction_id, tracker.announce_request.transaction_id);
	}
	response.interval = binary.BigEndian.Uint32((*buffer).Bytes()[8:]);
	response.leechers = binary.BigEndian.Uint32((*buffer).Bytes()[12:]);
	response.seeders = binary.BigEndian.Uint32((*buffer).Bytes()[16:]);

	var id = 0;
	var i int = 20;
	for i < ret {
		ip := net.IPv4((*buffer).Bytes()[i], (*buffer).Bytes()[i + 1], (*buffer).Bytes()[i + 2], (*buffer).Bytes()[i + 3]);
		// fmt.Println(binary.LittleEndian.Uint32(ip.To4()), "vs", ip.To4());
		port := binary.BigEndian.Uint16((*buffer).Bytes()[i + 4:]);
		fmt.Printf("[%d] %s:%d\n", id, ip.To4().String(), port);
		id++;
		i += 6;
	}
	if i != ret {
		// fmt.Printf("%d vs %d\n", i, ret);
		panic(fmt.Sprintf("%d vs %d\n", i, ret));
	} else {
		fmt.Println("OK!");
	}
	fmt.Println("Seeders: ", binary.BigEndian.Uint32((*buffer).Bytes()[16:]));
	fuckit(ret, *buffer);
	return nil;
}

func (tracker *Tracker) ipv4_write_announce(conn *net.UDPConn) (error) {
	buffer := make([]byte, 98);
	// var announcement ipv4_announce_request;

	tracker.announce_request.action = 1;
	tracker.announce_request.transaction_id = 42;
	copy(tracker.announce_request.peer_id[:20], "42");
	// tracker.announce_request.peer_id = []byte("42");
	tracker.announce_request.downloaded = 0;
	tracker.announce_request.left = 0;
	tracker.announce_request.uploaded = 0;
	tracker.announce_request.event = 2;
	tracker.announce_request.ip_addr = 0;
	tracker.announce_request.key = 0;
	tracker.announce_request.num_want = -1;
	// tracker.announce_request.num_want = int32(math.Floor(float64((packetsize - 20)) / float64(10)));
	tracker.announce_request.port = 6861;
	binary.BigEndian.PutUint64(buffer, tracker.announce_request.connection_id);
	binary.BigEndian.PutUint32(buffer[8:], tracker.announce_request.action);
	binary.BigEndian.PutUint32(buffer[12:], tracker.announce_request.transaction_id);
	copy(buffer[16:36], tracker.announce_request.info_hash[:]);
	copy(buffer[36:56], tracker.announce_request.peer_id[:]);
	binary.BigEndian.PutUint64(buffer[56:], tracker.announce_request.downloaded);
	binary.BigEndian.PutUint64(buffer[64:], tracker.announce_request.left);
	binary.BigEndian.PutUint64(buffer[72:], tracker.announce_request.uploaded);
	binary.BigEndian.PutUint32(buffer[80:], tracker.announce_request.event);
	binary.BigEndian.PutUint32(buffer[84:], tracker.announce_request.ip_addr);
	binary.BigEndian.PutUint32(buffer[88:], tracker.announce_request.key);
	binary.BigEndian.PutUint32(buffer[92:], uint32(tracker.announce_request.num_want));
	binary.BigEndian.PutUint16(buffer[96:], tracker.announce_request.port);
	// binary.BigEndian.PutUint32(buffer[16:], tracker.announce_request.info_hash);
	// binary.BigEndian.PutUint32(buffer[36:], tracker.announce_request.info_hash);

	if ret, err := conn.Write(buffer); err != nil || ret != 98 {
		if (err == nil) {
			fmt.Printf("Sent %d instead of 98\n", ret);
			return (errors.New("Didn't send the whole packet!"));
		}
		return (err);
	}
	fmt.Println("Sent packet!");
	// tracker.announce_request.downloaded
	// announcement.announce_request.action = 1;
	// announcement.announce_request.transaction_id = 42;
	// announcement.
	return nil;
}

func (tracker *Tracker) announcement(conn *net.UDPConn) (error) {
	var err	error;

	if err = tracker.ipv4_write_announce(conn);err != nil {
		return err;
	}
	if err = tracker.ipv4_read_announce(conn);err != nil {
		return err;
	}
	return nil;
}

