## Name
Auth Rest API service

## Description
This Service Below API for user Signup, Signin , Authenticate access token, Revoke access token and Refresh token

1. Sign up (creation of user) using email and password  - /signup   
    Sample Request 

   {   
    "email" :"test2904@gmail.com",
    "password": "test"
    }

2. Sign in - /login
     Authentication of user credentials
     JWT tokens (access token, refresh token) are returned in response 
     access token expire after 2 hrs , refresh after 7 days

    Sample Request 

   {   
    "email" :"test2904@gmail.com",
    "password": "test"
    }

3. Authorization of token  - /authorize-token 
     - Mechanism of sending token along with a request from client to service
     - checks for valid signature and expiry 
    sample request payload :
      {
    "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyOTA0QGdtYWlsLmNvbSIsImV4cCI6MTczMTQxOTM4OSwidXNlcklkIjoxfQ.D9Yk5JhnupeeiH52iRv2dyCsZ9mdr1O4nhHYT_yTAbQ"
     }
4. Revocation of token - /revoke-token 
      revoking a token 
      sample request payload:
      {
    "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyOTA0QGdtYWlsLmNvbSIsImV4cCI6MTczMTQxOTM4OSwidXNlcklkIjoxfQ.D9Yk5JhnupeeiH52iRv2dyCsZ9mdr1O4nhHYT_yTAbQ"
    }
5. Refresh a token  -  /refresh-token
    Renew Access token using refresh token
    sample request payload
    {
    "refreshtoken" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyOTA0QGdtYWlsLmNvbSIsImV4cCI6MTczMjAxNzQzMCwidXNlcklkIjoxfQ.aW0BM82DwcVx2Qht3BT9iuhyqqyLbqCRdtOhvN-ot0k"
}


## RUN 
To Run  Application go run cmd/main.go start

## DB 
Used file store DB "api.db" 

## tests
go test -v ./...


## Swagger Documentation

once the server is up , documentation is available at  http://localhost:8080/swagger/index.html#/


## CURLS

curl --location 'http://localhost:8080/signup' \
--header 'Content-Type: application/json' \
--data-raw '{   
    "email" :"test32904@gmail.com",
    "password": "test"
}'

curl --location 'http://localhost:8080/login' \
--header 'Content-Type: application/json' \
--data-raw '{   
    "email" :"test22904@gmail.com",
    "password": "test"
}'

curl --location 'http://localhost:8080/authorize-token' \
--header 'Authorization;' \
--header 'Content-Type: application/json' \
--data '{
    "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyMjkwNEBnbWFpbC5jb20iLCJleHAiOjE3MzI3ODk0NDksInVzZXJJZCI6M30.9o2Bkjp8XdzY5Dkjfb5gVnC_V09gLi9PFD3ZnFtk1Xs"
}'


curl --location 'http://localhost:8080/revoke-token' \
--header 'Content-Type: application/json' \
--data '{
    "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyMjkwNEBnbWFpbC5jb20iLCJleHAiOjE3MzIxOTEyMTIsInVzZXJJZCI6M30.NuZgO4lPqQx24Shl1S6jPovTo5ct4JsHe7X5tkvCUlg"
}'


curl --location 'http://localhost:8080/refresh-token' \
--header 'Content-Type: application/json' \
--data '{
    "refreshtoken" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyMjkwNEBnbWFpbC5jb20iLCJleHAiOjE3MzI3ODg4MTIsInVzZXJJZCI6M30.-P9Xdfrmt0vt3EKDUWdiINe98Zph-fcX9eRGYHWxmqs"
}'


## Postman Collection 

Kindly Refer to "AuthRestAPI.postman_collection.json" import them on your postman to test api's properly