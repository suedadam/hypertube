/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   scraper.go                                         :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: asyed <asyed@student.42.fr>                +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2018/05/29 13:25:58 by asyed             #+#    #+#             */
/*   Updated: 2018/05/29 14:43:18 by asyed            ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"net/url"
	"strings"
)

const piratebayurl string = "https://thepiratebay.org/search/%s/0/99/0"
/*
** Return magnet link.
*/

func searchTitle(title string) (string) {
	var	magnet_url string;

	resp, err := soup.Get(fmt.Sprintf(piratebayurl, url.PathEscape(title)));
	if err != nil {
		return "";
	}
	body := soup.HTMLParse(resp);
	links := body.Find("div", "id", "SearchResults").FindAll("a");
	for _, link := range links {
		if link.Pointer == nil {
			break ;
		}
		magnet_url = link.Attrs()["href"];
		if (strings.HasPrefix(magnet_url, "magnet:?")) {
			return magnet_url;
		}
	}
	return "";
}
