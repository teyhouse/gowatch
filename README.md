# 📁 GOWATCH - Keep track of changed files 
GOWATCH will keep track off changed files of your choice based on SHA256-Filesum.  
Just define files or folders (or mixed-mode) within the settings.json.
   
![screenshot](screenshot.png?raw=true)
  
  
Example settings.json:
```
{
    "files": {
      "test": "/home/administrator/test.txt",
      "dir" : "/home/administrator/www"
    }
}
```

[![goreleaser](https://github.com/teyhouse/gowatch/actions/workflows/main.yml/badge.svg)](https://github.com/teyhouse/gowatch/actions/workflows/main.yml) 
    
# 📃 Requirements
Always make sure to have your settings.json in the same folder as the gowatch-binary.  
The hashes.json will be recreated, if not moved or deleted. The file event.json is also optional.
  
  
# 🛠️ Usage
If you want to check for file changes on a regular basis, you should add GOWATCH to your Crontab.  
For example, to check every minute for changed files:  
```* * * * * /usr/bin/gowatch/gowatch```
  
It would also be possible to check in realtime, using watch:
```watch /home/administrator/gowatch/gowatch```
  
# 📝 Logs
Every time GOWATCH detects a change in one of your files, it will be logged in your syslog:  
![screenshot](log.png?raw=true)
  
You can also use the file event.json to trigger a custom HTTP-GET, which will automatically attach a new GET-Parameter to your URI, called ``&logmessage``, containing the base64-encoded log-message.  
  
# 💭 Debugmode
You can use the ```--debug``` flag / arg to get more details on which changes GOWATCH has detected:  
![screenshot](debug.png?raw=true)
  
## © License
MIT