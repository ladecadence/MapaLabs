# MapaLabs

Web app for the Map of the Encuentro de Labs

## Installation

The backend is programmed using Go, after you install Go SDK (https://go.dev), download this repo and you can compile the app with the Go build system:
```
$ cd MapaLabs
$ go mod tidy
$ go build cmd/mapalabs.go
```

After this you'll have a "mapalabs" binary in the current directory. You can run this locally with a configuration file (example provided).

```
$ ./mapalabs -c config.toml
```

The app uses a SQLite database, on the first run it will do a migration and create the tables for you. You'll need at least to add a user to be able to upload new entries. You can do this with the sqlite command line tool:
```
$ sqlite3 mapalabs.db

sqlite> insert into users (name, password, email, role) values ("username", "3368a8bb88fc4931eb655fb8ae35175c78bcf878f195f5ebffe0779951bbf309", "user@email.xyz", 0);

```
The password is sha256 hashed, you can generate a password on the command line:
```
$ echo -n "MySecurePassword" | sha256sum
```

## Deployment

You'll need a web server with proxy capabilities to redirect requests to the web app. Create a virtualhost for example, and make a proxy redirect to the port the app will be running (defined in the configuration file). For example if using NGINX web server:

```
server {
    root /var/www/mapalabs.example.net;
    index index.html index.htm index.nginx-debian.html;
    server_name mapalabs.ladecadence.net;
    location / {
                # Proxy pass and CORS
                proxy_pass http://localhost:8080;
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
                add_header Access-Control-Allow-Origin "*";
        }
    listen 443 ssl;
    
    # ssl certificates
    [...]

}
```

Then you can run the app as a service when the server starts. For example creating a SystemD service file at /etc/systemd/system/ (or the directory holding systemd services in your server):
```
Description=MapaLabs web service
After=network.target

[Service]
User=www-data
Group=www-data
ExecStart=/var/www/mapalabs.example.net/mapalabs -c /etc/mapalabs/config.toml
[Install]
WantedBy=multi-user.target
```

Then edit the configuration file and fill the url and main_path fields with the correct values.

And run the service after placing all the files in the correct directories (for example config in /etc/mapalabs/, mapalabs binary, mapalabs.db, html/ and static/ in /var/www/mapalabs.example.net, etc)

```
$ sudo systemd start mapalabs.service
```

If the web server is working and doing the proxy redirect, everything should be working.
