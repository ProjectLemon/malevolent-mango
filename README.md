# Malicious Mango
This is a portfolio service that hosts portfolios in pdf-format. If you couldn't
tell from the title.

## Usage
1. Setup go if not done so. Guide can be found in the official
   [documentation](https://golang.org/doc/install)
2. Download the project using:

    ```
   go get github.com/ProjectLemon/malicious-mango
   ```
   or simply clone the repository into your go source folder.

2. If you download the repo using git clone you will have to install the
   dependencies separatley using:
   ```
    go get
  ```

3. Install the project by, from inside the project folder, running 
   ```
   go install
   ```

4. If you have $GOPATH/bin in your path you can just type *malicious-mango* to
   start the server. Otherwise, start the server using:
   ```
   $GOPATH/bin/malicious-mango
   ```

5. Open your browser and navigate to *localhost:8080* to use the server

##Database
This project asumes the usage of a database (although it will build even without one). To set up this simply provide a file named *.db_cnf* in the same folder that this file is in. **Remember to add this file to .gitignore if you're using a different one from this repo.** Inside this file you are to provide, username, password, drivername and datasourcename as in the following example:

```
username Alice
password SuperStrongPassword1234
drivername mysql
datasourcename tcp(0.0.0.0:3306)/databasename?parseTime=true
```

The file has to be present when the server is started. 
