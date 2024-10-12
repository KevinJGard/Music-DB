# MusicDB  
## Author: Kevin Jonathan Garduño Escobar.  
MusicDB is a music manager that connects to a database where it stores the user's music through a graphical interface, made in ___GO___.  

## Requirements
You need to have ___Go___ installed, version 1.23.1. You can install it from its official website. [Golang](https://go.dev/), you will also need to install ___sqlite3___.

## Installation
To configure and run the project locally, follow these steps:
1. Clone the repository:
	```bash
	git clone https://github.com/KevinJGard/MusicDB
	```
2. Navigate to the project directory:
	```bash
	cd MusicDB/
	```
3. Use the following command to start unloading the dependencies
	```bash
	go mod tidy
	```
4. If you want to run it you can do the following command:
	```bash
	go run src/interface.go
	```
5. If you want to compile it you can do the following command:
	```bash
	go build -o <name for the executable> src/interface.go
	```
6. To run the executable:
	```bash
	./<name for the executable>
	```
It may take some time for the interface to be displayed.  

## Use of the interface:  
The ___Miner___ menu contains two options  
* Set path  
This option is to be able to choose your directory with music.  
* Mine metadata  
This option starts the mining of mp3 files and shows the progress bar.  

The ___Options___ menu contains two options  
* Settings  
This option opens a new window with two buttons to switch between dark and light themes.  
* Help
This option opens the project's Github browser.  

The ___Screen___ menu contains two options  
* Full screen  
This option sets the window to full screen size.  
* Quit  
This option closes the program.  

After mining you will see a list of all your songs in the database, when you select one from the list on the right side of the screen you will have three buttons  
* Edit P.  
This button opens a new window for editing the performer, it will give you three options  
    * Person  
When you press this button, the entries to put the person's data are enabled, also the option to put him/her in a band is enabled.
    * Group  
When you click on it, the entries are enabled to put the data of the group.  
    * Undefined  
When you press this button you can only change the name of the performer.  
* Edit A.  
This button opens a new window where you can enter the new fields for the album to be modified.  
* Edit Song  
Opens a new window for editing the song data, where you can enter new data.  

When you click on ___Cancel___, the window closes and when you enable the button ___Submit___ the performer is modified and a notification is sent with your input.  
To see the changes you have to select another song and go back to the one you modified and you should see the changes reflected.  

### To make a search:  
You need to search according to the language set: 
- ar:\<Artist name\>&&\<Another artist\>&&\<Another artist\>  
- al:\<Album name\>&&\<Another name\>&&\<Another name\>&&\<Another name\>  
- ti:\<Song title\>&&\<Another title\>&&\<Another title\>&&\<Another title\>  
- ye:\<Year of song\>&&\<Another year\>&&\<Another year\>  
- ge:\<Genre\>&&\<Another genre\>&&\<Another genre\>  

To search multiple fields at once, use `||` to separate prefixes.  
For example:  
`
ti:Exist&&One Last Kiss||ar:Michael Jackson&&Coldplay&&Jose Jose||al:Thriller||ye:2018&&1986||ge:Pop  
`

To do your search you must press the button that is on the right side with the magnifying glass icon, this will open another window where the results of your search will be displayed.  

## Running Unit Test
To run the unit tests, navigate to the src directory and use the following command:  
```bash
go test ./test/
```

## Command Line Usage  
For testing or mining music metadata directly from the command line, you can use the main.go file It allows you to mine metadata and search songs without the graphical interface.  

Example usage:  
```bash
go run src/main.go <directory_path> <search_query>
```
Example:  
```bash
go run src/main.go /home/user/Music/ "ti:Exist||ar:Michael Jackson"
```
This will display the search results and MP3 metadata found in the given directory.

## Contributing
If you find any issues or have suggestions for improvements, feel free to open an issue or a pull request on the project's GitHub repository.