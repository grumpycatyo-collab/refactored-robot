# User Service API Gateway (REST API)

This service allows user registration, login, image upload and retrieval, getting user details, and token refresh.

## HTTP Routes:

* __POST /user__ : Allows user registration.
* __GET /user/get/:id__ : Retrives User information.
* __POST /user/login__ : Authenticates user login.
* __POST /user/setimg/:id__ : Allows image uploads for user.
* __GET /user/getimg/:id__ : Retrieve image for corresponding user.
* __POST /user/refresh__ : Refresh user's access token.

## gRPC Endpoints

* __AddUser__ : To register a new user.
* __GetUser__ : To retrieve user information.
* __DeleteUser__ : To delete user.
* __Login__ : To authenticate user login.
* __UploadImage__ : To upload image for user.
* __GetImage__ : To retrieve user image.
* __RefreshToken__ : To refresh user's access token.

# User Service API Gateway (REST API)

## HTTP Routes and their payloads

### User Registration (POST /user)
```
{
  "name": "John Doe",
  "password": "password123",
  "image": "<base64-encoded-image>"
}
```

### User Login (POST /user/login)
```
{
   "name": "John Doe",
   "password": "password123"
}
```

### Upload User Image (POST /user/setimg/{id})
The Post body should be form-data format with field.
The id is included in the URL.

Form-Data:
```
{
  "file": "<files>"
}
```
### Refresh Access Token (POST /user/refresh)
```
{
  "refreshToken": "<refresh-token-from-login-operation>"
}
```
Note: The image field is a base64-encoded image.

## GET Requests

### Get User Details (GET /user/get/{id})
Here the id of the user is passed as a path parameter in URL. Example:

```
/user/get/1
```
### Get User Image (GET /user/getimg/{id})
Here the id of the user is passed as a path parameter in URL. Example:

```
/user/getimg/1
```
Note for both of these requests, you need to pass in the token obtained from login operation via the Authorization header as Bearer token.

This is an initial version of the service. For any additional concerns, request, or bug reports, please create an issue or raise a pull request.

Note: The password encryption and decryption technique is not mentioned yet. It will be updated soon.