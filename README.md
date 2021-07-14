<h1 align="center">
  <br>
<img src="static/bcaem.png" width="200px" alt="bcaem">
</h1>

<h4 align="center">Fast AEM scope gathering tool for all your public and private BugCrowd Programs</h4>


<p align="center">
<a href="https://goreportcard.com/report/github.com/ethicalhackingplayground/bcaem"><img src="https://goreportcard.com/badge/github.com/ethicalhackingplayground/bcaem"></a>
<a href="https://github.com/ethicalhackingplayground/bcaem/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat"></a>
<a href="https://github.com/ethicalhackingplayground/bcaem/releases"><img src="https://img.shields.io/github/release/ethicalhackingplayground/bcaem"></a>
<a href="https://twitter.com/z0idsec"><img src="https://img.shields.io/twitter/follow/z0idsec.svg?logo=twitter"></a>
<a href="https://discord.gg/MQWCem5b"><img src="https://img.shields.io/discord/862900124740616192.svg?logo=discord"></a>
</p>

<p align="center">
  <a href="#install">Install</a> •
  <a href="#usage">Usage</a> •
  <a href="#examples">Examples</a> •
  <a href="https://discord.gg/MQWCem5b">Join Discord</a> 
</p>

---

If you love to test for Adobe Experience Manager Vulnerabilities well, I'm proud to announce BCAEM (Bugcrowd AEM) is a tool that will gathering
all Adobe Experience Manager programs from your public or private programs.


## Installation
Make sure you've a recent version of the Go compiler installed on your system.
Then just run:
```
GO111MODULE=on go get -u github.com/ethicalhackingplayground/bcaem
```

## Usage
```
▶ bcaem bc -t <session-token> <other-flags>
```
How to get the session token:
- Bugcrowd: login, then grab the `_crowdcontrol_session` cookie

Remember that you can use the --help flag to get a description for all flags.

## Examples
Below you'll find some example commands.

### Print all in-scope AEM targets from bugcrowd
```
▶ bcaem bc -t <YOUR_TOKEN> -b 
```
The output will look like this:
```
app.example.com
*.user.example.com
*.demo.com
www.something.com
```

### Print all in-scope aem targets from all your private Bugcrowd programs that offer rewards
```
▶ bcaem bc -t <YOUR_TOKEN> -b -p
```

### License

Erebus is distributed under [GPL-3.0 License](https://github.com/ethicalhackingplayground/erebus/blob/main/LICENSE)

<h1 align="left">
  <a href="https://discord.gg/MQWCem5b"><img src="static/Join-Discord.png" width="380" alt="Join Discord"></a>
</h1>