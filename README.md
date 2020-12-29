# LeakyTap  
  
At first this seems pretty pointless, but actually it might help someone else? My first paid bounty was based around identifying leaked PHP source. The server wasn't interpreting/rendering it and just sending it back to the client.

This tool just tries to help automate that by allowing a user to pipe a list of URLs in via stdin and spitting out any results.

## Installing  
```
go get github.com/crawl3r/LeakyTap
```  
  
## Standard Run  
```
cat urls.txt | ./leakytap
```
  
## Run and save the output to file  
```
cat urls.txt | ./leakytap -o output.txt
```  
  
## Run in quiet mode, only prints the identified leaky endpoints
```
cat urls.txt | ./leakytap -q
```
  
