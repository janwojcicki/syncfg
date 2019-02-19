# syncfg
Convenient config files synchronizer.

## Why?
Imagine you want to control your config files with git. It is certainly possible and, if you care about good change tracking probably better than this solution, but it is quite cumbersome. A simple change in one of the config files is at least 3 commands (add, commit, push). With syncfg it is shortened to only one call

## How?
This repository includes both a simple web server for storing files, as well as a program that is actually responsible for making changes to your machine. The files are transported via http. The password/registration system is currently VERY unsafe so I would not recommend storing anything of value here

## How do I use it?
### Setup
If you just want to try this app I recommend you try it on my server. The setup is relatively simple.

```
git clone https://github.com/janwojcicki/syncfg.git
cd syncfg/app
go build syncfg.go
```

and now you have the syncfg app built. If you want to make it globally accessible either modify PATH or just:
```
sudo mv syncfg /usr/bin
```

After installation the app requires you to register yourself to some server
```
syncfg register [server-host] [username] [password]
```

Again, this app is VERY unsafe. It stores your password in plain text in the config files, both on your pc and on the server. In case of ANY break in your password WILL BE vulnerable. Please, use something that isn't your email password

### Usage
#### Add file
```
syncfg add [repository-name] [file-name] [file-path]
```
Repository will be automatically created if you don't have one with this name

#### Commit changes
```
syncfg commit
```

#### Get your files
```
syncfg get [repository-name]
```

## Contributing
The main flaw I see in this program is the way it handles passwords and registrations. It should very much be fixed.
