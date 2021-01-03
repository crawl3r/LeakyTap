# LeakyTap  
  
At first this seems pretty pointless, but actually it might help someone else? My first paid bounty was based around identifying leaked PHP source. The server wasn't interpreting/rendering it and just sending it back to the client.

This tool just tries to help automate that by allowing a user to pipe a list of URLs in via stdin and spitting out any results.

Any issues, let me know! I plan to extend the source code it can identify, so if there are specific requests please let me know.

### Thanks  
  
Big thanks to Hakluke, I used Hakrawler's (https://github.com/hakluke/hakrawler) concurrency and picked at the concurrency/goroutine code to patch mine.

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
  
### License  
I'm just a simple skid. Licensing isn't a big issue to me, I post things that I find helpful online in the hope that others can:  
 A) learn from the code  
 B) find use with the code or   
 C) need to just have a laugh at something to make themselves feel better  
  
Either way, if this helped you - cool :)  
