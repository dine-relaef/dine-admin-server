
## Run
To start the application, use the following command:
```sh
docker-compose up --build
```

## Commit
To set up the commit script, run the following command once:
```sh
chmod +x commit.sh
```
For subsequent commits, use the script with your commit message:
```sh
./commit.sh "your message" branch_name
```
