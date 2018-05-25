/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   handshake.go                                       :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: asyed <asyed@student.42.fr>                +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2018/05/22 19:16:07 by asyed             #+#    #+#             */
/*   Updated: 2018/05/25 10:07:37 by asyed            ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"fmt"
	"net"
	"encoding/binary"
	"bytes"
	"errors"
	// "reflect"
)

type udp_tracker_response struct {
	action			uint32;
	transaction_id	uint32;
	connection_id	uint64; //Can remove this.
}

type udp_tracker_request struct {
	protocol_id		uint64;
	action			uint32;
	transaction_id	uint32;
}

func (tracker *Tracker) read_handshake(transaction_id uint32, conn *net.UDPConn) (error) {
	var response	udp_tracker_response;
	var buffer		*bytes.Buffer;

	if buffer = bytes.NewBuffer(make([]byte, 16)); buffer == nil {
		return errors.New("Failed to initialize buffer");
	}
	if ret, err := conn.Read((*buffer).Bytes()); err != nil || ret != 16 {
		if (err == nil) {
			fmt.Printf("Read %d instead of 16\n", ret);
			return errors.New("Didn't read the whole packet!");
		}
		return err;
	}
	response.action = binary.BigEndian.Uint32((*buffer).Bytes());
	if response.transaction_id = binary.BigEndian.Uint32((*buffer).Bytes()[4:8]); response.transaction_id != transaction_id {
		return errors.New("Invalid transaction id");
	}
	tracker.announce_request.connection_id = binary.BigEndian.Uint64((*buffer).Bytes()[8:16]);
	fuckit(16, *buffer);
	return nil;
}

func (tracker *Tracker) write_handshake(transaction_id uint32, conn *net.UDPConn) (error) {
	buffer := make([]byte, 16);

	binary.BigEndian.PutUint64(buffer, 0x41727101980);
	binary.BigEndian.PutUint32(buffer[8:], 0);
	binary.BigEndian.PutUint32(buffer[12:], transaction_id);
	if ret, err := conn.Write(buffer); err != nil || ret != 16 {
		if (err == nil) {
			fmt.Printf("Sent: %d instead of 16\n", ret);
			return (errors.New("Didn't send the whole packet!"));
		}
		return (err);
	}
	return (nil);
}

func (tracker *Tracker) handshake(conn *net.UDPConn) (error) {
	var buffer			bytes.Buffer;
	var transaction_id	uint64;
	var err				error;

	buffer.Grow(16);
	transaction_id = 42;
	if err = tracker.write_handshake(uint32(transaction_id), conn); err != nil {
		return err;
	}
	if err = tracker.read_handshake(uint32(transaction_id), conn); err != nil {
		return err;
	}
	return nil;
}
