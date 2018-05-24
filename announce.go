/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   announce.go                                        :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: asyed <asyed@student.42.fr>                +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2018/05/23 12:53:54 by asyed             #+#    #+#             */
/*   Updated: 2018/05/24 12:58:19 by asyed            ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"fmt"
	"errors"
)

func (tracker *Tracker) ipv4_announce_request() (error) {
	// var announcement ipv4_announce_request;

	tracker.announce_request.action = 1;
	tracker.announce_request.transaction_id = 42;
	tracker.announce_request.peer_id = []byte(42);
	tracker.announce_request.downloaded = 0;
	tracker.announce_request.left = 0;
	tracker.announce_request.uploaded = 0;
	tracker.announce_request.event = 2;
	tracker.announce_request.ip_addr = 0;
	tracker.announce_request.port = 6861;
	// tracker.announce_request.downloaded
	// announcement.announce_request.action = 1;
	// announcement.announce_request.transaction_id = 42;
	// announcement.
	return nil;
}

func (tracker *Tracker) announcement() (error) {
	var err	error;

	if err = tracker.ipv4_announce_request();err != nil {
		return err;
	}
	return nil;
}

