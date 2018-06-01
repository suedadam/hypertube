package main

import (
		"encoding/csv"
		"fmt"
		"os"
		"strconv"
)

type Movie struct {
		Id		string;
		Title	string;
		Genres	string;
		Rating	float64;
		NumVotes int;
}

// func convMovieToCSVRecord(m Movie) []string{
// 	ret := []string{m.Id, m.Title,genreStr,strconv.Itoa(m.Rating)};
// 	return ret;
// }

func main(){
	 // read data from CSV file
	 csvFile, err := os.Open(os.Args[1]);

	 if err != nil {
			 fmt.Println(err);
	 }

	 
	 reader := csv.NewReader(csvFile);
	 reader.LazyQuotes = true;
	 fmt.Println("Setting up new reader to be ready to read file for basic movie metadata");
	 reader.FieldsPerRecord = -1;
	 reader.Comma = '\t';
	 
	 var oneMovie Movie;
	 var allMovies map[string]*Movie;
	 allMovies = make(map[string]*Movie);
	 totalMovieRecords := 0;
	 for err == nil {
		 oneRecord, loc_err := reader.Read();
		 if len(oneRecord) == 0 {
			 fmt.Println(len(oneRecord));
			 fmt.Println(loc_err);
			 // os.Exit(1)
			}
			//  fmt.Println("Read one record");
			if(len(oneRecord) > 0 && oneRecord[1]  == "movie"){
				 oneMovie.Id = oneRecord[0];
				 // fmt.Println(oneRecord[2]);
				 oneMovie.Title = oneRecord[2];
				 oneMovie.Genres = oneRecord[8];
				 oneMovie.Rating = -1;
				 oneMovie.NumVotes = -1;
				 oneMovieCopy := oneMovie;
				 allMovies[oneMovie.Id] = &oneMovieCopy;
				 totalMovieRecords++;
				}
				err = loc_err; 
			}
			fmt.Println("Read file");
			fmt.Println("Total Records ", totalMovieRecords);
			csvFile.Close();
			
			// read data from CSV file
			csvFile, err = os.Open(os.Args[2]);
			
			if err != nil {
				fmt.Println(err);
			}
			
			defer csvFile.Close();
			
			reader2 := csv.NewReader(csvFile);
			reader2.LazyQuotes = true;
			fmt.Println("Setting up new reader2 to be ready to read file for ratings");
			reader2.FieldsPerRecord = -1;
			reader2.Comma = '\t';
			for err == nil {
				 oneRecord, loc_err := reader2.Read();
				 if len(oneRecord) == 0 {
					 fmt.Println(len(oneRecord));
					 fmt.Println(loc_err);
						 // os.Exit(1)
					  }
					  if(len(oneRecord) > 0){
						  val, ok := allMovies[oneRecord[0]];
						  if(ok){
							  fmt.Println("getting ratings");
							  val.Rating, err = strconv.ParseFloat(oneRecord[1], 64);
							  val.NumVotes, err = strconv.Atoi(oneRecord[2]);
							 }
							}
							err = loc_err; 
						}
					  //  fmt.Println("Read one record");
			// sanity check
			// NOTE : You can stream the JSON data to http service as well instead of saving to file
	//fmt.Println(string(jsondata))
	
	// now write to JSON file
	outputFile := "movies.csv";
	csvOutputFile, err := os.Create(outputFile);
	defer csvOutputFile.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Created output csv file");
	csvWriter := csv.NewWriter(csvOutputFile);
	defer csvWriter.Flush();
	totalMovieRecords = 0;
	for i, value := range allMovies {
		if(totalMovieRecords % 1000 == 0){
			fmt.Println("Write record ", i);
		}
		err := csvWriter.Write([]string{value.Id, value.Title,value.Genres,fmt.Sprint(value.Rating),strconv.Itoa(value.NumVotes)});
		if err != nil {
			fmt.Println(err)
		}
		totalMovieRecords++;
	}
	fmt.Println("Wrote to file.");
}