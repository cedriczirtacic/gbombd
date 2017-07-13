# gbombd
ZIP 💣  httpd. Sends a ZIP bomb if it recognizes the client as a bot or a scanner. You're going to find User Agents inside the code so you can add more to those.
It simply creates a **big** file with repetitive data (in this case zeros), so the compresion is more effective but when the client decompresses the gzip'ed data it will consume all his CPU rendering it useless.

Got the idea from [here](https://blog.haschek.at/2017/how-to-defend-your-website-with-zip-bombs.html).
## Compile
Build it using the Go complier (you don't need to *go get* anything):
```bash
$ go build -v -o gbombd gbombd.go
```

## Using this thing
```bash
$ ./gbombd -help
Usage of ./gbombd:
  -filename string
    	Filename to create and use (default "bomb.gz")
  -port int
    	HTTPd port to use (default 80)
  -preserve
    	Preserve the bomb file for future use
  -size int
    	Size of file to create (default 10240)
  -verbose
    	Be verbose
```

By default it will open port 80 and you will need more privileges for that. The default file size is 10G and name **bomb.gz**. The file will be removed once the application exits but you can avoid that by using the **-preserve** switch.

## Why? Go is fun and annoying some people is even funnier.

