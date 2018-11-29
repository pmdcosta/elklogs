# Elklogs

`elklogs` is a command line utility to query and tail ELK (elasticsearch, logstash, kibana) logs from the terminal.

`elklogs` is a fork of the original https://github.com/knes1/elktail project with some core changes to the project structure and features.

Even though it's powerful, using Kibana's web interface to search and analyse the logs is not always practical.
Sometimes you just wish to tail -f the logs that you normally view in kibana to see what's happening right now.
Elklogs allows you to do just that, and more. Tail the logs. Search for errors and specific events on commandline.
Pipe the search results to any of the standard unix tools. Use it in scripts.
Redirect the output to a file to effectively download a log from es / kibana etc...

## Status

WIP: The current implementation is in its first stages.
