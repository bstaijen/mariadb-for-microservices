FROM golang:1.7

# Download and install any required third party dependencies into the container.
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/gorilla/mux
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/urfave/negroni
RUN go get github.com/bstaijen/mariadb-for-microservices/shared/util
RUN go get github.com/bstaijen/mariadb-for-microservices/shared/models
RUN go get github.com/bstaijen/mariadb-for-microservices/shared/helper
RUN go get github.com/joho/godotenv
RUN go get github.com/Sirupsen/logrus
RUN go get github.com/meatballhat/negroni-logrus
RUN go get gopkg.in/DATA-DOG/go-sqlmock.v1

# 
ADD . /go/src/mariadb.com/profile-service/
WORKDIR /go/src/mariadb.com/profile-service
RUN go build main.go

# Expose port 5000 to the host so we can access our application
EXPOSE 5000

# Tell Docker what command to run when the container starts
CMD ["/go/src/mariadb.com/profile-service/main"]