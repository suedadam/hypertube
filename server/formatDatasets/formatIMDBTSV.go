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
		Rating	int;
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

	 defer csvFile.Close();

	 reader := csv.NewReader(csvFile);
	 fmt.Println("Setting up new reader to be ready to read file");
	 reader.FieldsPerRecord = -1;
	 reader.Comma = '\t';
	 
	 var oneMovie Movie;
	 var allMovies []Movie;
	 totalMovieRecords := 0;
	 for err == nil {
		 oneRecord, loc_err := reader.Read();
		 //  if loc_err != nil {
			 // 	 fmt.Println(len(oneRecord));
			 // 	 fmt.Println(loc_err);
			 // 	 os.Exit(1)
			 //  }
			 fmt.Println("Read one record");
			 if(len(oneRecord) > 0 && oneRecord[1] == "movie"){
				 oneMovie.Id = oneRecord[0];
				 oneMovie.Title = oneRecord[2];
				 oneMovie.Genres = oneRecord[8];
				oneMovie.Rating = 0;
				allMovies = append(allMovies, oneMovie);
				totalMovieRecords++;
			}
			err = loc_err; 
		}
		fmt.Println("Read file");
		fmt.Println("Total Records ", totalMovieRecords);
		
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
	for i, value := range allMovies {
		if(i % 1000 == 0){
			fmt.Println("Write record ", i);
		}
		err := csvWriter.Write([]string{value.Id, value.Title,value.Genres,strconv.Itoa(value.Rating)});
		if err != nil {
			fmt.Println(err)
		}
		
	}
	fmt.Println("Wrote to file.");
}