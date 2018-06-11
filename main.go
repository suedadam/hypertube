/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   main.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: asyed <asyed@student.42.fr>                +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2018/05/22 16:54:47 by asyed             #+#    #+#             */
/*   Updated: 2018/06/01 19:07:41 by asyed            ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"fmt"
	// "net"
	// "bytes"
	"time"
	"os"
	"strings"
	"net/http"
	"net/url"
	"errors"
	// "C"
	// "encoding/hex"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
)

var StatsTable		map[string]*Torrent_info;
var DefaultClient	torrent.Config;

type Torrent_info struct {
	c *torrent.Client;
	t *torrent.Torrent;
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodPost) {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed);
		return ;
	}
	if err := r.ParseForm();err != nil {
		http.Error(w, "", http.StatusNoContent);
		return ;
	}
	for key, val := range r.PostForm {
		if strings.Compare(key, "title") == 0 {
			fmt.Println("Title = ", val);
			if magnet := searchTitle(val[0]);magnet != "" {
				fmt.Fprint(w, magnet);
			} else {
				http.NotFound(w, r);
				return ;
			}
		}
	}
	http.Error(w, "500", http.StatusInternalServerError);
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	var tinfo Torrent_info;

	if (r.Method != http.MethodPost) {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed);
		return ;
	}
	if err := r.ParseForm();err != nil {
		http.Error(w, "", http.StatusNoContent);
		return ;
	}
	for key, val := range r.PostForm {
		if strings.Compare(key, "magnet") == 0 {
			// fmt.Println("Magnet = ", val);
			if displaypath, err := tinfo.download_torrent(val[0]);err == nil {
			// if magnet := searchTitle(val[0]);magnet != "" {
				fmt.Fprint(w, fmt.Sprintf("http://localhost:1338/%s", url.PathEscape(displaypath)))
				// fmt.Fprint(w, "LOL IT WORKED?!");
				return ;
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError);
				return ;
			}
		}
	}
	http.Error(w, "500", http.StatusInternalServerError);
}

/*
** What do we want to do if there's multiple video files? Grab the largest?
** ToDo: Conversion.
*/

func (tinfo *Torrent_info) download_torrent(magnet_url string) (string, error) {
	var largest *torrent.File;

	c, err := torrent.NewClient(&DefaultClient);
	if err != nil {
		return "", err;
	}
	t, err := c.AddMagnet(magnet_url);
	if err != nil {
		return "", err;
	}
	fmt.Println("Waiting for info data.");
	<-t.GotInfo();
	fmt.Println("Got info data.");
	files := t.Files();
	largest = nil;
	for i := 0; i < len(files); i++ {
		if !strings.HasSuffix(files[i].DisplayPath(), ".mkv") && strings.HasSuffix(files[i].DisplayPath(), ".mp4") {
			files[i].SetPriority(torrent.PiecePriorityNone);
			continue ;
		}
		if (largest == nil || files[i].Length() >= largest.Length()) {
			if (largest != nil) {
				files[i].SetPriority(torrent.PiecePriorityNone);
			}
			largest = files[i];
		}
		// files[i] //Make a better way to check for file types.
	}
	if largest == nil {
		return "", errors.New("No supported file types.");
	}
	t.DownloadAll();
	tinfo.c = c;
	tinfo.t = t;
	go func() {
		for ;; {
			if (t.BytesCompleted() == t.Info().TotalLength()) {
				return ;
			}
			fmt.Printf("%s -> [%d - %d]\n", largest.DisplayPath(), t.BytesCompleted(), t.Info().TotalLength())
			c.WriteStatus(os.Stdout);
			time.Sleep(time.Second);
		}
	}()
	// first := largest.firstPieceIndex.Int();
	// last := largest.firstPieceIndex.Int();
	// largest.t.updatePiecePriorities(first, ((last - first) / 100) * 5);
	return largest.DisplayPath(), nil;
}

func init() {
	var baseDir storage.ClientImpl;

	baseDir = storage.NewFile("storage");
	DefaultClient.DefaultStorage = baseDir;
}

func main() {
	mux := http.NewServeMux();
	mux.HandleFunc("/search", searchHandler);
	mux.HandleFunc("/download", downloadHandler);
	// DefaultClient = NewFile("storage");
	go func() {
		panic(http.ListenAndServe(":1337", mux));
	}()
	go func() {
		// http.Handle("/", http.FileServer(http.Dir("storage")));
		panic(http.ListenAndServe(":1338", http.FileServer(http.Dir("storage"))));
	}()
	fmt.Println("Listening!!!");
	select{}
	// const sample string = "magnet:?xt=urn:btih:bee75372b98077bfd4de8ef03eb33e9289be5cd8&dn=Avengers+Infinity+War+2018+NEW+PROPER+720p+HD-CAM+X264+HQ-CPG&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969";
	// // const sample string = "magnet:?xt=urn:btih:3f2a441e2c4b84f25e44403328aeffc432a15ae7&dn=Ready+Player+One.2018.HDRip.X264.AC3-EVO%5BN1C%5D&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969";
	// var tinfo Torrent_info;

	// go func() {
	// 	if err := tinfo.download_torrent(sample); err != nil {
	// 		panic(err);
	// 	}
	// 	for {
	// 		info := tinfo.t.Info();
	// 		if info == nil {
	// 			fmt.Println("Info == nil")
	// 		} else {
	// 			fmt.Printf("%d vs %d\n", tinfo.t.BytesCompleted(), info.TotalLength());
	// 			time.Sleep(time.Second);
	// 		}
	// 		stats := tinfo.t.Stats();
	// 		fmt.Printf("Total = %d Pending = %d Active = %d Connected = %d Half = %d\n", stats.TotalPeers, stats.PendingPeers, stats.ActivePeers, stats.ConnectedSeeders, stats.HalfOpenPeers);
	// 		if (tinfo.t.BytesCompleted() == info.TotalLength()) {
	// 			break ;
	// 		}
	// 	}
	// 	tinfo.c.Close();
	// }()
	// select{}
}