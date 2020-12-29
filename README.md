# LeakyTap  
  
At first this seems pretty pointless, but actually it might help someone else? My first paid bounty was based around identifying leaked PHP source. The server wasn't interpreting/rendering it and just sending it back to the client.

This tool just tries to help automate that by allowing a user to pipe a list of URLs in via stdin and spitting out any results:
  
```
cat urls.txt | ./leakytap
```
  
```
cat urls.txt | ./leakytap -o output.txt
```  
  
```
cat urls.txt | ./leakytap -q
```

